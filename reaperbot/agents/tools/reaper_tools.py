import asyncio

from pydantic_ai import ModelRetry
from dotenv import load_dotenv
from devtools import debug
from agents.tools.reaper_api import (ReaperAPI,
                              Domain,
                              ScannedDomain,
                              GetScannedDomainsResult,
                              ScanDomainResult,
                              GetLiveHostsResult,
                              GetLiveEndpointsResult,
                              GetAttackEndpointResult,
                              GetRequestsForEndpointResult)
from utils.logging import send_log_message

load_dotenv()

async def reaper_get_scanned_domain_by_id(domain_id: int) -> ScannedDomain:
    """
    Retrieves the record about a domain that reaper has already scanned.
    Use this tool to avoid rescanning a domain that has already been scanned.
    """
    try:
        await send_log_message(f"Tool: GetScannedDomainById input: {domain_id}")
        result = await ReaperAPI().get_scanned_domain_by_id(domain_id=domain_id)
        await send_log_message(f"Tool: GetScannedDomainById result: {result}")
        return result
    except:
        raise ModelRetry(f'Could not get the domain_id: {domain_id} metadata')

async def reaper_get_scanned_domains() -> GetScannedDomainsResult:
    """
    Retrieves a list of records about domains that reaper has already scanned.
    Use this tool to avoid rescanning a domain that has already been scanned.
    """
    try:
        await send_log_message(f"Tool: GetScannedDomains input: -")
        result = await ReaperAPI().get_scanned_domains()
        await send_log_message(f"Tool: GetScannedDomains result: {result}")
        return result
    except:
        raise ModelRetry('Could not get the domains list')

async def reaper_scan_domain(domain: Domain) -> ScanDomainResult:
    """
    Passively scan a domain to trigger reaper to collect/know about possible hosts
    in that domain. This tool does NOT retrieve them.  It only kicks off a scan
    for them and waits until the scan is complete before the tool returns.  It
    does not return any results except for the domain_id and a status or error 
    message. Live hosts are retrieved only by the get_live_hosts_for_domains tool.

    Be sure to get the domain to scan from the user and prompt them for one
    instead of making one up.  Do not use example.com, for instance.

    Args:
        domain: str - The domain to start a subdomain/host finding scan
    """
    try:
        await send_log_message(f"Tool: ScanDomain input: {domain.name}")
        result = await ReaperAPI().scan_domain(domain.name)
        await send_log_message(f"Tool: ScanDomain result: {result}")
        return result
    except:
        raise ModelRetry('Could not get the live hosts')


async def reaper_get_live_hosts_for_domains(domain_id: int, domains: list[str]) -> GetLiveHostsResult:
    """
    Get the list of live hosts from a passive scan tool named reaper. This
    should be called after a domain has been scanned by the scan_domain
    tool (to know the domain_id). Call this tool as often as needed to get
    the most up to date information.

    Args:
        domains: list[str] - The list of domains to filter live hosts by
    """
    if domains == []:
        raise ModelRetry('Provide the domain or domains to get live hosts for in the domains list[str]')
    try:
        await send_log_message(f"Tool: GetLiveHosts input: {domain_id} {domains}")
        result = await ReaperAPI().get_live_hosts(domain_id=domain_id, domains=domains)
        await send_log_message(f"Tool: GetLiveHosts result: {result}")
        return result
    except:
        raise ModelRetry('Could not get the live hosts')


async def reaper_get_live_endpoints_for_domains(domains: list[str], methods: list[str], filter_on_params_usage: bool) -> GetLiveEndpointsResult:
    """
    Get the list of valid and live endpoints from a passive scan tool named reaper. This
    must be called after the get_live_hosts_for_domains tool has been made so that reaper
    has knowledge of the live endpoints associated with those hosts. Call this tool
    as often as needed to get the most up to date information.

    Set filter_on_params_usage to False by default. Set filter_on_params_usage
    to True to filter endpoints that allow for GET requests with non-empty url parameters 
    or POST requests with keys in the body params tend to be candidates for BOLA or IDOR attacks.

    Args:
        domains: A list of domains to filter the list of endpoints on. If an empty list, no domain filter is applied
        methods: A list of http methods e.g. ["GET", "POST"] to filter the list of endpoints. If an empty list, return all live endpoints for all methods
        filter_on_params_usage: Whether or not to filter endpoints based on params usage
    """
    try:
        await send_log_message(f"Tool: GetLiveEndpoints input: {domains} {methods} {filter_on_params_usage}")
        result = await ReaperAPI().get_live_endpoints(domains=domains, methods=methods, params=filter_on_params_usage)
        endpoints = []
        for endpoint in result.endpoints:
            endpoints.append(f"{endpoint.endpoint_id} {endpoint.method} {endpoint.host}{endpoint.path} params: {endpoint.params}, ")
        await send_log_message(f"Tool: GetLiveEndpoints result: {endpoints}")
        return result
    except:
        raise ModelRetry('Could not get the live endpoints')


async def reaper_get_requests_for_endpoint_id(endpoint_id: int) -> GetRequestsForEndpointResult:
    """
    For a given endpoint (by endpoint_id), fetch requests that reaper has seen to that endpoint.
    Use that context when determining how to proceed with an attack before
    carrying one out. For example, to know which param_key would be best to
    attack from an attacker's perspective.

    Call this tool before running an attack to have the proper context and
    example requests with full payloads first.

    Args:
        endpoint_id: the reaper endpoint_id (int) to fetch full request and responses for to get proper context on parameters and body structure and actual example values
    """
    try:
        await send_log_message(f"Tool: GetRequestsForEndpoint endpoint_id: {endpoint_id}")
        result = await ReaperAPI().get_requests_for_endpoint(endpoint_id=endpoint_id)
        reqs = []
        for req in result.requests:
            reqs.append(f"{req.id} {req.method} {req.url} params: \"{req.param_keys}\", body: \"{req.body_keys}\", ")
        await send_log_message(f"Tool: GetRequestsForEndpoint result: {reqs}")
        return result
    except:
        raise ModelRetry('Could not fetch requests for that endpoint')


async def reaper_test_attack_endpoint_id(endpoint_id: int, param_key: str) -> GetAttackEndpointResult:
    """
    Send an endpoints' endpoint_id to be used to carry out a parameter fuzzing
    attack for BOLA/IDOR using the passive scan tool named reaper. Always use 
    the get_requests_for_endpoint_id tool first to determine the one parameter 
    key that is most likely to be most dangerous for users and most beneficial 
    from an attacker's perspective if successful. Only use param_key values
    from returned endpoints data.

    Args:
        endpoint_id: the reaper endpoint_id (int) to probe
        param_key: the key name of the parameter to focus the fuzzing test on.
    """
    try:
        await send_log_message(f"Tool: AttackEndpoint input: {endpoint_id} {param_key}")
        result = await ReaperAPI().get_attack_endpoint(endpoint_id=endpoint_id, param_key=param_key)
        await send_log_message(f"Tool: AttackEndpoint result: {result}")
        return result
    except:
        raise ModelRetry('Could not attack that endpoint')


async def reaper_get_attack_results(attack_id: int) -> str:
    """
    Fetch attack results by attack_id for tests/attacks already run using the passive scan tool named
    reaper. Ensure the test_attack_endpoint_id tool has been run successfully
    first.

    Args:
        attack_id: the attack_id to fetch the results for in JSON
    """
    try:
        await send_log_message(f"Tool: GetAttackResults input: {attack_id}")
        result = await ReaperAPI().get_attack_results(attack_id=attack_id, limit=25)
        await send_log_message(f"Tool: GetAttackResults result: {result.status} {result.message}")
        return result
    except:
        raise ModelRetry('Could not fetch attack results for that attack_id')