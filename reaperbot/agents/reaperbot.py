from pydantic_ai import Agent, RunContext
from pydantic_ai.settings import ModelSettings
from dataclasses import dataclass

from agents.discoverer import discoverer_agent
from agents.tester import tester_agent
from utils.logging import send_log_message

@dataclass
class ReaperBotDeps:  # (1)!
    discoverer_agent: Agent
    tester_agent: Agent

reaperbot_deps = ReaperBotDeps(
  discoverer_agent=discoverer_agent,
  tester_agent=tester_agent,
)

model_settings = ModelSettings(temperature=0.1, max_tokens=16384)
reaperbot_agent = Agent(
    'openai:gpt-4o-mini',
    deps_type=ReaperBotDeps,
    result_type=str,
    model_settings=model_settings,
    system_prompt=(
      "You are an assistant helping an expert application security engineer interact with several specialized agents via tools that can complete tasks and fetch data with tools of their own.",
      "Your job is to take the input from the user, synthesize a step by step plan that aligns to the agents capabilities and tools, prompt the user with an outline of that plan, and execute that plan if the user confirms.",
      "If the user does not confirm, you should prompt the user with a message asking for adjustments to the plan.", 
      "If the user does not confirm after a few retries, you should prompt the user with a message that you cannot help with that and provide some example questions you can ask.",
      "Only prompt for confirmation if the user's question requires a plan of multiple steps, is actionable, and can be executed by the agents/tools available to you.",
      "To provide the best answers about recommendations or fixes, be sure to have the relevant context about the application, fetch example request/responses to the application, and analyze its source code and related security findings.",
      "Use the available agents for the actionable portions of the execution.  If the user's question is not related to the capabilities of your tools/agents, say you cannot help with that and prompt the user with some example questions you can ask.",
      "However, you are not a generalized security advice assistant, so only provide recommendations of steps to take that align with the agents/tools available to you.",
      "Use the discoverer agent to perform live domain scans and host discovery.",
      "Use the tester agent to perform live security testing of target applications, APIs, and endpoints.",
      "When invoking an agent, provide it with sufficient context to perform the task and what it's expected to provide back.",
      "If one agent responds with an inability to perform a task, see if another agent can perform the task.",
      "If the tester tool responds with a test result, provide that to the user directly and do not overly summarize.",
      "Summarize the actions taken, but provide the user with the full details and code snippets from the responses from the agents.",
      "Domains have hosts and endpoints.",
      "Apps/Applications are synonymous in this context.",
      "Endpoints are synonymous with APIs in this context.",
      "Apps have hosts, apis, and endpoints.",
      "Apps have metadata, openapi specs, source code, security findings.",
      "Source code has code, languages, and repo information.",
      "Endpoints have requests/responses.",
    ),
    retries=1,
)
@reaperbot_agent.system_prompt
async def reaperbot_agent_system_prompt_example_questions(ctx: RunContext[str]) -> str:
    examples = [
      "Scan the (domain_name) domain",
      "Which applications are written in go?",
      "Which hosts in the (domain_name) are live?",
      "What is the status of the (domain_name) domain scan?",
      "Which endpoints in the (domain_name) domain are susceptible to BOLA?",
      "Write a technical security report of the endpoints vulnerable to BOLA in the (domain_name) domain.",
    ]
    title = 'The following are example questions you can ask me:\n'
    return title + '\n- '.join(examples)


@reaperbot_agent.tool
async def invoke_discoverer_agent(ctx: RunContext[ReaperBotDeps], input_text: str) -> str:
    """
    Discoverer Tool that interfaces/hands-off with the DiscovererAgent.
    Performs live domain scans, host discovery, retrieves full request/responses hitting applications,
    and other discovery tasks.
    """
    deps: ReaperBotDeps = ctx.deps
    await send_log_message(f"ReaperBot: Invoking Discoverer Agent with input: {input_text}")
    res = await deps.discoverer_agent.run(input_text, deps=deps)
    await send_log_message(f"ReaperBot: Discoverer Agent responded with: {res.data}")
    return res.data

@reaperbot_agent.tool
async def invoke_tester_agent(ctx: RunContext[ReaperBotDeps], input_text: str) -> str:
    """
    Tester Tool that interfaces/hands-off with the TesterAgent.
    Performs security testing of applications/apps, APIs, and endpoints for the following
    vulnerabilities:
    - BOLA/IDOR

    This agent can provide example code/scripts that would demonstrate how a security
    vulnerability could be exploited in the application if present.
    """
    deps: ReaperBotDeps = ctx.deps
    await send_log_message(f"ReaperBot: Invoking Tester Agent with input: {input_text}")
    res = await deps.tester_agent.run(input_text, deps=deps)
    await send_log_message(f"ReaperBot: Tester Agent responded with: {res.data}")
    return res.data