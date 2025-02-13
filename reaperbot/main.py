import uvicorn
from typing import List
from dotenv import load_dotenv
from datetime import datetime
from asyncio import sleep
from devtools import debug
from contextlib import asynccontextmanager
from fastapi import FastAPI, WebSocket, Request, Depends
from fastapi.templating import Jinja2Templates
from fastapi.responses import HTMLResponse, FileResponse
from fastapi.staticfiles import StaticFiles
from starlette.websockets import WebSocketDisconnect

from utils.database import Database, get_db, get_ws_db
from agents.reaperbot import reaperbot_agent, reaperbot_deps
from utils.logging import send_log_message
from utils.websockets import log_websockets

load_dotenv()

@asynccontextmanager
async def lifespan(_app: FastAPI):
    async with Database.connect() as db:
        yield {'db': db}

app = FastAPI(lifespan=lifespan)

templates = Jinja2Templates(directory="templates")
app.mount("/static", StaticFiles(directory="static"), name="static")

chat_responses: List[str] = []

@app.websocket("/log")
async def log_websocket_endpoint(websocket: WebSocket):
    """WebSocket endpoint for logging."""
    await websocket.accept()
    log_websockets.append(websocket)
    try:
        while True:
            await websocket.receive_text()
    except WebSocketDisconnect:
        if websocket in log_websockets:
            log_websockets.remove(websocket)

@app.websocket("/ws")
async def websocket_endpoint(websocket: WebSocket, db: Database = Depends(get_ws_db)):
    """WebSocket endpoint for real-time AI responses."""
    await websocket.accept()
    await sleep(0.2)  # Add a small delay to ensure the connection is established
    await send_log_message("ReaperBot is ready to assist.")
    try:
        while True:
            user_message = await websocket.receive_text()
            await send_log_message(f"ReaperBot: Processing user request: {user_message}")
            messages = await db.get_messages("desc")

            async with reaperbot_agent.run_stream(user_message, deps=reaperbot_deps, message_history=messages) as result:
                # debug(reaperbot_agent)
                async for message in result.stream_text(delta=True):
                    await websocket.send_text(message)
            await db.add_messages(result.new_messages_json())
            await send_log_message("ReaperBot: Response complete.")
    except WebSocketDisconnect:
        await send_log_message("ReaperBot offline.")
    except Exception as e:
        await send_log_message(f"Error: {str(e)}")
        await send_log_message("ReaperBot offline.")

@app.post("/clear_messages")
async def clear_messages(request: Request, db: Database = Depends(get_db)):
    """Endpoint to clear the chat messages in the database."""
    await db.clear_messages()
    return {"status": "success"}

@app.get("/favicon.ico", include_in_schema=False)
async def favicon():
    """Serve the favicon.ico file."""
    return FileResponse("static/favicon.ico")

@app.get("/", response_class=HTMLResponse)
async def chat_page(request: Request):
    """Serve the chat page."""
    return templates.TemplateResponse("index.html", {"request": request, "chat_responses": chat_responses})

if __name__ == "__main__":
    uvicorn.run("main:app", host="0.0.0.0", port=11000, log_level="info", reload=True)
