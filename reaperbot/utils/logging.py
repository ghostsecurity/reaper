from datetime import datetime
from starlette.websockets import WebSocketDisconnect
from .websockets import log_websockets

async def send_log_message(message: str):
    """Send a log message to all connected log WebSocket clients."""
    timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    log_message = f"[{timestamp}] {message}"
    for websocket in log_websockets:
        try:
            await websocket.send_text(log_message)
        except WebSocketDisconnect:
            if websocket in log_websockets:
                log_websockets.remove(websocket)
