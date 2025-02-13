import asyncio
from contextlib import asynccontextmanager
from functools import partial
from pathlib import Path
import sqlite3
from dataclasses import dataclass
from typing import Any, AsyncIterator, Callable, LiteralString, TypeVar, List
from concurrent.futures.thread import ThreadPoolExecutor
from typing_extensions import Annotated, ParamSpec
from pydantic import Field, TypeAdapter
from pydantic_ai.messages import ModelMessage, ModelMessagesTypeAdapter
from datetime import datetime
from fastapi import WebSocket, Request

THIS_DIR = Path(__file__).parent

MessageTypeAdapter = TypeAdapter(
    Annotated[ModelMessage, Field(discriminator='kind')]
)
P = ParamSpec('P')
R = TypeVar('R')


@dataclass
class Database:
    """Rudimentary database to store chat messages in SQLite.

    The SQLite standard library package is synchronous, so we
    use a thread pool executor to run queries asynchronously.
    """

    con: sqlite3.Connection
    _loop: asyncio.AbstractEventLoop
    _executor: ThreadPoolExecutor

    @classmethod
    @asynccontextmanager
    async def connect(
        cls, file: Path = THIS_DIR / '.chat_app_messages.sqlite'
    ) -> AsyncIterator['Database']:
        loop = asyncio.get_event_loop()
        executor = ThreadPoolExecutor(max_workers=1)
        con = await loop.run_in_executor(executor, cls._connect, file)
        slf = cls(con, loop, executor)
        try:
            yield slf
        finally:
            await slf._asyncify(con.close)

    @staticmethod
    def _connect(file: Path) -> sqlite3.Connection:
        con = sqlite3.connect(str(file))
        cur = con.cursor()
        cur.execute(
            'CREATE TABLE IF NOT EXISTS messages (id INT PRIMARY KEY, message_list TEXT);'
        )
        con.commit()
        return con

    async def add_messages(self, messages: bytes):
        await self._asyncify(
            self._execute,
            'INSERT INTO messages (id, message_list) VALUES (?, ?);',
            int(datetime.now().timestamp()),
            messages,
            commit=True,
        )
        await self._asyncify(self.con.commit)

    async def get_messages(self, dir: str) -> List[ModelMessage]: # type: ignore
        if not dir:
            dir = "desc"
        c = await self._asyncify(
            self._execute, f"SELECT message_list FROM messages order by id {dir}"
        )
        rows = await self._asyncify(c.fetchall)
        messages: List[ModelMessage] = [] # type: ignore
        for row in rows:
            messages.extend(ModelMessagesTypeAdapter.validate_json(row[0]))
        return messages
    
    async def clear_messages(self):
        await self._asyncify(
            self._execute, "DELETE FROM messages"
        )
        await self._asyncify(self.con.commit)

    def _execute(
        self, sql: LiteralString, *args: Any, commit: bool = False
    ) -> sqlite3.Cursor:
        cur = self.con.cursor()
        cur.execute(sql, args)
        if commit:
            self.con.commit()
        return cur

    async def _asyncify(
        self, func: Callable[P, R], *args: P.args, **kwargs: P.kwargs
    ) -> R:
        return await self._loop.run_in_executor(  # type: ignore
            self._executor,
            partial(func, **kwargs),
            *args,  # type: ignore
        )

async def get_ws_db(ws: WebSocket) -> Database:
    return ws.state.db

async def get_db(request: Request) -> Database:
    return request.state.db
