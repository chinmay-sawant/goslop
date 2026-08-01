from fastapi import FastAPI
import time

app = FastAPI()
STORE = {}


@app.get("/x")
def read_x(k: str):
    global STORE
    STORE[k] = 1
    return STORE


@app.get("/slow")
async def slow():
    time.sleep(1)
    return {"ok": True}
