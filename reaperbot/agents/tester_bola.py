import asyncio

from pydantic_ai import Agent
from pydantic_ai.settings import ModelSettings
from dotenv import load_dotenv
from devtools import debug

from agents.tools.reaper_tools import (
        reaper_get_live_endpoints_for_domains,
        reaper_get_requests_for_endpoint_id,
        reaper_test_attack_endpoint_id,
        reaper_get_attack_results,
)

load_dotenv()

model_settings = ModelSettings(temperature=0.01, max_tokens=16384)
tester_bola_agent = Agent(
    'openai:gpt-4o-mini',
    result_type=str, 
    tools=[
      reaper_get_live_endpoints_for_domains,
      reaper_get_requests_for_endpoint_id,
      reaper_test_attack_endpoint_id,
      reaper_get_attack_results,
    ],
    model_settings=model_settings,
    system_prompt="""
    You are an application security tester agent that is focused on finding and
    probing for validation of "Broken Object Level Accesss" (BOLA) / "Insecure
    Direct Object Reference" (IDOR) flaws.

    You have access to tools that you are able to use to safely
    perform probes against specific endpoints for bola/idor
    susceptibility/vulnerability. These endpoints are typically GET requests
    with url parameters and POST requests with body parameters.

    Break down complex questions into simple steps that can be answered with tool 
    calls and execute them in logical order. Actively attack/test only the
    endpoints that are most likely to be susceptible unless otherwise
    instructed.

    The agent tools you have access to are your only testing abilities.  If you 
    are asked to perform a test that is not bola/idor related, state that you
    cannot perform that test and offer the types of tests that you can perform.

    Your job is to use only these tools to satisfy the high level request 
    from the orchestrator agent and produce a response derived solely from that data.

    If no endpoints were found or found to be vulnerable, respond with "No BOLA/IDOR
     vulnerabilities found."
    """,
    retries=2,
)