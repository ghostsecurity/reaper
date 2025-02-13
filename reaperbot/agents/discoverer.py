import asyncio

from pydantic_ai import Agent
from dotenv import load_dotenv
from devtools import debug
from pydantic_ai.settings import ModelSettings
from agents.tools.reaper_tools import (
        reaper_get_scanned_domain_by_id,
        reaper_get_scanned_domains,
        reaper_scan_domain,
        reaper_get_live_hosts_for_domains,
        reaper_get_live_endpoints_for_domains,
        reaper_get_requests_for_endpoint_id,
)

load_dotenv()

model_settings = ModelSettings(temperature=0.01, max_tokens=16384)
discoverer_agent = Agent(
    'openai:gpt-4o-mini',
    result_type=str, 
    tools=[
      reaper_get_scanned_domain_by_id,
      reaper_get_scanned_domains,
      reaper_scan_domain,
      reaper_get_live_hosts_for_domains,
      reaper_get_live_endpoints_for_domains,
      reaper_get_requests_for_endpoint_id,
    ],
    model_settings=model_settings,
    system_prompt="""
      You are an advanced tool calling assistant interacting on behalf of
      an expert application security engineer with a system called
      Reaper to perform web application enumeration and testing actions to
      gather data about target domains, hosts, and endpoints.

      The goal is to help the engineer perform some or all of the following 
      tasks for a given domain:
      1. Use a tool (reaper) to scan a domain to discover live hosts if not already
      scanned by reaper.
      2. Use a tool (reaper) to retrieve live hosts in that domain
      3. Use a tool (reaper) to retrive live endpoints on those hosts with certain
      characteristics useful for planning a test attack.

      The primary starting point for a workflow is a domain name, and most
      workflows will be focused on that domain until the user changes to another
      domain to focus on.  Break down complex questions into simple steps that 
      can be answered with tool calls and execute them in logical order.

      Important: If you haven't been given a domain to focus on by
      the user, prompt the user to choose one. Do not make up a domain to scan.
      Do not use example.com, for example.  If a tool is not available, do not 
      make up responses.

      You are encouraged to run these tools and use their response details to decide
      which step to take next and run more tools to get the information you
      need to answer the request.
    """,
    retries=2,
)