import os
import re
import json
import httpx
import asyncio
import time
from dotenv import load_dotenv
from pydantic import BaseModel, Field
from typing import Literal, Optional
from tenacity import retry, stop_after_attempt, wait_exponential, RetryError


class Domain(BaseModel):
    name: str = Field(..., description="The domain name to be scanned by Reaper", example="ghostbank.net")

class ScannedDomain(BaseModel):
    id: int
    name: str
    status: str

class GetScannedDomainsResult(BaseModel):
    domains: list[ScannedDomain] = []

class ScanDomainResult(BaseModel):
    status: Literal["success", "error"]
    domain_id: int = 0
    content: str

class GetLiveHostsResult(BaseModel):
    status: Literal["success", "error"]
    hosts: list[str] = []

class LiveEndpoint(BaseModel):
    endpoint_id: int
    method: Literal["GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"]
    host: str
    path: str
    params: list[str]

class GetLiveEndpointsResult(BaseModel):
    endpoints: list[LiveEndpoint]

class EndpointRequestResponse(BaseModel):
    body: str
    content_length: int
    content_type: str
    created_at: str
    headers: str
    id: int
    request_id: int
    status: str
    status_code: int
    updated_at: str

class EndpointRequestRequest(BaseModel):
    body: str
    body_keys: str
    content_length: int
    content_type: str
    created_at: str
    header_keys: str
    headers: str
    host: str
    id: int
    method: str
    param_keys: str
    proto: str
    proto_major: int
    proto_minor: int
    response: EndpointRequestResponse
    source: str
    url: str

class GetRequestsForEndpointResult(BaseModel):
    requests: list[EndpointRequestRequest]


class GetAttackEndpointResult(BaseModel):
    attack_id: int
    endpoint_id: int
    param_key: str
    status: Literal["success", "completed", "error"]
    message: str

class AttackResult(BaseModel):
    created_at: str
    endpoint: str
    fuzz_attack_id: int
    hostname: str
    id: int
    ip_address: str
    port: int
    request: str
    response: str
    scheme: Literal["http","https"]
    status_code: int
    updated_at: str
    url: str

class GetAttackResults(BaseModel):
    attack_results: list[AttackResult]
    status: Literal["success", "error"]
    message: str


class ReaperAPI:
    def __init__(self):
        load_dotenv()
        self.api_key = os.environ['X_REAPER_TOKEN']
        self.base_url = os.environ['REAPER_BASE_URL']
        self.headers = {
            "Accept": "application/json",
            "X-Reaper-Token": self.api_key
        }

    class ScanStatusError(Exception):
        """Custom exception to trigger retries."""
        pass

    async def get_scanned_domain_by_id(self, domain_id: int) -> ScannedDomain:
        """Get the metadata of an already scanned domain by its id"""

        self.url = f"{self.base_url}/api/scan/domains/{domain_id}"

        async with httpx.AsyncClient() as client:
            response = await client.get(self.url, headers=self.headers)

        domain = {}
        if response.status_code == 200:
            data = response.json()

            domain["id"] = data.get('id')
            domain["name"] = data.get('name')
            domain["status"] = data.get('status')
            return ScannedDomain(**domain)
        else:
            raise


    async def get_scanned_domains(self) -> GetScannedDomainsResult:
        """Get the list of already scanned domains"""

        self.url = f"{self.base_url}/api/scan/domains"

        async with httpx.AsyncClient() as client:
            response = await client.get(self.url, headers=self.headers)

        domains = []
        if response.status_code == 200:
            data = response.json()

            for res in data:
                domain = {}
                domain["id"] = res.get('id')
                domain["name"] = res.get('name')
                domain["status"] = data.get('status')
                domains.append(domain)

        return GetScannedDomainsResult(domains=domains)


    async def scan_domain(self, domain: str) -> ScanDomainResult:
        """Start scan for the domain and handle the status."""
        if not domain:
            return ScanDomainResult(status="error", content="Domain is required but not provided.")

        self.url = f"{self.base_url}/api/scan/domains"

        # Start the scan
        start_response = await self._post_scan_request(domain)
        if start_response.status == "error":
            return start_response

        # Wait for scan completion or handle conflict
        try:
            scan_result = await self._wait_for_scan_completion(domain)
            return scan_result
        except RetryError:
            return ScanDomainResult(status="error", content="The scan was unable to complete in time after multiple retries.")


    async def get_live_hosts(self, domain_id: int, domains: list[str]) -> GetLiveHostsResult:
        if domain_id == 0:
            return []

        if domains is None or domains == []:
            domains = ['.*']
        else:
            domains = [f".*\\.{domain}|{domain}" for domain in domains]

        self.url = f"{self.base_url}/api/scan/domains/{domain_id}/hosts"
        async with httpx.AsyncClient() as client:
            response = await client.get(self.url, headers=self.headers)

        livehosts = []

        if response.status_code == 200:
            data = response.json()
            for host in data:
                if host.get('status') == "live":
                    host_domain = host.get('name', '')
                    if any(re.match(pattern, host_domain) for pattern in domains):
                        livehosts.append(host.get('name'))

        return GetLiveHostsResult(status="success", hosts=livehosts)


    async def get_live_endpoints(self, domains: list[str], methods: list[str], params: bool = True) -> GetLiveEndpointsResult:
        if domains is None or domains == []:
            domains = ['.*']
        else:
            domains = [f".*\\.{domain}|{domain}" for domain in domains]

        if methods is None or methods == []:
            methods = ["GET", "POST", "PUT", "PATCH", "OPTIONS", "DELETE"]

        self.url = f"{self.base_url}/api/endpoints"

        async with httpx.AsyncClient() as client:
            response = await client.get(self.url, headers=self.headers)

        endpoints = []

        if response.status_code == 200:
            data = response.json()
            for host in data:
                if host.get('method').upper() in methods:
                    host_domain = host.get('hostname', '')
                    if any(re.match(pattern, host_domain) for pattern in domains):
                        if params == False:
                            # Don't filter by presence of params
                            endpoints.append(host)
                        else:
                            # only add ones with get url params or post body params
                            if host.get("params") != "":
                                endpoints.append(host)

        live_endpoints = []
        for ep in endpoints:
            entry = {}
            entry["endpoint_id"] = ep.get('id')
            entry["method"] = ep.get('method')
            entry["host"] = ep.get('hostname')
            entry["path"] = ep.get('path')
            entry["params"] = ep.get('params').split(",")
            live_endpoints.append(entry)

        return GetLiveEndpointsResult(endpoints=live_endpoints)


    async def get_requests_for_endpoint(self, endpoint_id: int) -> GetRequestsForEndpointResult:
        if endpoint_id == 0:
            return []

        self.url = f"{self.base_url}/api/endpoints/{endpoint_id}"
        async with httpx.AsyncClient() as client:
            response = await client.get(self.url, headers=self.headers)

        ep_id = 0
        ep_hostname = ""
        ep_method = ""
        ep_path = ""
        if response.status_code == 200:
            ep = response.json()
            if ep.get('id'):
                ep_id = ep.get('id')
                ep_hostname = ep.get('hostname')
                ep_method = ep.get('method')
                ep_path = ep.get('path')
        if ep_id == 0:
            return json.dumps([])

        reqs = []
        self.url = f"{self.base_url}/api/requests"
        async with httpx.AsyncClient() as client:
            response = await client.get(self.url, headers=self.headers)

        if response.status_code == 200:
            pattern = f".*{ep_path}$"
            data = response.json()
            for req in data:
                if req.get('method').upper() == ep_method.upper() and req.get('host').lower() == ep_hostname.lower() and re.match(pattern, req.get('url')):
                    reqs.append(req)

        return GetRequestsForEndpointResult(requests=reqs)


    async def get_attack_endpoint(self, endpoint_id: int, param_key: str) -> GetAttackEndpointResult:
        if endpoint_id is None:
            return "endpoint_id required but not provided"
        if param_key is None or param_key == "":
            return "param_keys required but not provided"

        self.url = f"{self.base_url}/api/attack"

        data = {}
        data["endpoint_id"] = int(endpoint_id)
        data["params"] = [param_key]
        async with httpx.AsyncClient() as client:
            response = await client.post(self.url, headers=self.headers, json=data)

        if response.status_code == 200:
            tries = 0
            max_tries = 8
            status = response.json().get('status')
            attack_id = response.json().get('id')
            while tries < max_tries:
                self.url = f"{self.base_url}/api/attacks/{attack_id}"
                async with httpx.AsyncClient() as client:
                    response = await client.get(self.url, headers=self.headers)

                if response.status_code == 200:
                    status = response.json().get('status')
                    result = GetAttackEndpointResult(
                            attack_id=int(attack_id),
                            endpoint_id=int(endpoint_id),
                            param_key=str(param_key),
                            status="success",
                            message="")

                if status == "success":
                    result.message = "Attack completed and found a vulnerable endpoint."
                    return result
                elif status == "completed":
                    result.status = "completed"
                    result.message = "Attack completed but did not find a vulnerable endpoint. Try again with another param_key."
                    return result
                else:
                    tries = tries + 1
                    time.sleep(2 * 1.2 * tries)
        else:
            return GetAttackEndpointResult(
                attack_id=0,
                endpoint_id=int(endpoint_id),
                param_key=str(param_key),
                status="error",
                message=f"There was an error launching the attack {response}")


    async def get_attack_results(self, attack_id: int, limit: int) -> GetAttackResults:
        if limit is None or limit == 0:
            limit = 50
        if attack_id is None or attack_id == 0:
            return "attack_id required but not provided"

        self.url = f"{self.base_url}/api/attacks/{attack_id}/results?limit={limit}"

        async with httpx.AsyncClient() as client:
            response = await client.get(self.url, headers=self.headers)

        results = []
        if response.status_code == 200:
            data = response.json()
            for res in data:
                results.append(res)
            
            return GetAttackResults(attack_results=results, status="success",
                                    message=f"Successfully fetched the attack results")
        else:
            return GetAttackResults(attack_results=results, status="error",
                                    message="Error fetching the attack results")


    @retry(stop=stop_after_attempt(3), wait=wait_exponential(min=2, multiplier=1.2))
    async def _post_scan_request(self, domain: str) -> ScanDomainResult:
        """Send a request to start scanning the domain."""
        data = {"name": domain, "auto_scan": True}
        async with httpx.AsyncClient() as client:
            response = await client.post(self.url, headers=self.headers, json=data)
        if response.status_code == 409:
            return ScanDomainResult(status="success", content=f"Scan already completed for domain '{domain}'.")
        if response.status_code == 200:
            domain_id = response.json().get('id')
            return ScanDomainResult(status="success", domain_id=domain_id, content=f"Scan started for domain '{domain}'.")
        else:
            return ScanDomainResult(status="error", content=f"Error starting scan for domain '{domain}': {response.text}")


    async def _get_domain_info(self, domain: str) -> Optional[dict]:
        """Fetch domain scan info from the server."""
        async with httpx.AsyncClient() as client:
            response = await client.get(self.url, headers=self.headers)
        if response.status_code != 200:
            return None
        domains = response.json()
        return next((dom for dom in domains if dom.get('name') == domain), None)


    @retry(stop=stop_after_attempt(10), wait=wait_exponential(min=2, multiplier=1.1))
    async def _wait_for_scan_completion(self, domain: str) -> ScanDomainResult:
        """Wait and check the status of a scan with retries."""
        domain_info = await self._get_domain_info(domain)
        if domain_info:
            status = domain_info.get('status')
            domain_id = domain_info.get('id')

            if status == "completed":
                return ScanDomainResult(status="success", domain_id=domain_id, content=f"Domain scan of: {domain} completed successfully.")
            else:
                raise self.ScanStatusError(f"Scan is in {status} status. Retrying...")

        return ScanDomainResult(status="error", content=f"Error fetching scan status for domain '{domain}'.")

    async def _handle_conflict(self, domain: str) -> ScanDomainResult:
        """Handle cases where the scan is already in progress or completed."""
        domain_info = await self._get_domain_info(domain)
        if domain_info:
            status = domain_info.get('status')
            domain_id = domain_info.get('id')

            if status == "completed":
                return ScanDomainResult(status="success", domain_id=domain_id, content=f"Domain scan has already been performed.")

        return ScanDomainResult(status="error", content=f"Unable to fetch domain status or domain '{domain}' does not exist.")