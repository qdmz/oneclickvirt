import json, urllib.request, urllib.error

BASE = "http://localhost:8888/api/v1"

def req(method, path, token=None, data=None):
    url = BASE + path
    body = json.dumps(data).encode() if data is not None else None
    r = urllib.request.Request(url, data=body, method=method)
    if token:
        r.add_header("Authorization", "Bearer " + token)
    if body is not None:
        r.add_header("Content-Type", "application/json")
    try:
        with urllib.request.urlopen(r, timeout=10) as resp:
            return resp.status, json.loads(resp.read().decode() or "{}")
    except urllib.error.HTTPError as e:
        try:
            payload = json.loads(e.read().decode() or "{}")
        except Exception:
            payload = {}
        return e.code, payload

def login(username, password):
    st, js = req("POST", "/auth/login", data={"username": username, "password": password})
    tok = js.get("data", {}).get("token") or js.get("data", {}).get("access_token")
    print(f"[login] {username} -> {st} token={'yes' if tok else 'no'}")
    return tok

admin = login("admin", "admin123456")
user = login("testuser", "test123456")

print("\n=== admin stats (before) ===")
st, js = req("GET", "/admin/tickets/stats", admin)
print(st, json.dumps(js.get("data")))

print("\n=== user creates ticket ===")
st, js = req("POST", "/user/tickets", user, {"title": "E2E测试工单", "content": "这是一个自动化测试工单内容。", "priority": "high"})
print(st, json.dumps(js))
tid = js.get("data", {}).get("id")
print("ticket id =", tid)

print("\n=== user lists tickets ===")
st, js = req("GET", "/user/tickets?page=1&pageSize=10", user)
print(st, "total=", js.get("data", {}).get("total"), "list=", [t["id"] for t in js.get("data", {}).get("list", [])])

print("\n=== admin lists tickets ===")
st, js = req("GET", "/admin/tickets?page=1&pageSize=10", admin)
print(st, "total=", js.get("data", {}).get("total"))

print("\n=== admin replies to ticket ===")
st, js = req("POST", f"/admin/tickets/{tid}/replies", admin, {"content": "您好，我们已经收到您的问题，正在处理。"})
print(st, json.dumps(js))

print("\n=== user replies back ===")
st, js = req("POST", f"/user/tickets/{tid}/replies", user, {"content": "谢谢，期待尽快解决。"})
print(st, json.dumps(js))

print("\n=== admin ticket detail ===")
st, js = req("GET", f"/admin/tickets/{tid}", admin)
d = js.get("data", {})
print(st, "status=", d.get("status"), "replyCount=", d.get("replyCount"), "replies=", len(d.get("replies", [])))

print("\n=== admin update status -> answered ===")
st, js = req("PUT", f"/admin/tickets/{tid}/status", admin, {"status": "answered"})
print(st, json.dumps(js))

print("\n=== admin stats (after) ===")
st, js = req("GET", "/admin/tickets/stats", admin)
print(st, json.dumps(js.get("data")))

print("\n=== admin closes ticket ===")
st, js = req("POST", f"/admin/tickets/{tid}/close", admin)
print(st, json.dumps(js))

print("\n=== admin deletes ticket ===")
st, js = req("DELETE", f"/admin/tickets/{tid}", admin)
print(st, json.dumps(js))
print("\nDONE")
