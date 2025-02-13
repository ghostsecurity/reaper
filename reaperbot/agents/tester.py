import asyncio

from pydantic_ai import Agent
from pydantic_ai.settings import ModelSettings
from dotenv import load_dotenv

from agents.tester_bola import tester_bola_agent
from utils.logging import send_log_message

load_dotenv()

model_settings = ModelSettings(temperature=0.01, max_tokens=16384)
tester_agent = Agent(
    'openai:gpt-4o-mini',
    result_type=str, 
    model_settings=model_settings,
    system_prompt="""
    You are an application security tester agent that is focused on finding and
    probing for validation of specific application security flaws.

    You have access to your own specialized agents via tools that you are
    to delegate the actual testing actions to.  Do not perform tests on your own. The
    agent tools you have access to are your only testing abilities. If you are
    asked to perform a test that you do not have a tool for, state that you
    cannot perform that test and offer the types of tests that you can perform.

    Your job is to use only these tools to satisfy the high level request 
    from the orchestrator agent and produce a response derived solely from that data.
    """,
    retries=2,
)

@tester_agent.tool_plain
async def invoke_tester_bola_agent(input_text: str) -> str:
    """Tester Tool that interfaces/hands-off with the BOLA specific TesterAgent."""
    await send_log_message(f"Tester BOLA Agent: invoke_tester_bola_agent input: {input_text}")
    res = await tester_bola_agent.run(input_text)
    return res.data
