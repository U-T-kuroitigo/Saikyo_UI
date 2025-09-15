// フロントJS：フォーム送信で /api/menu を叩いて結果表示
const form = document.getElementById("menuForm");
const statusEl = document.getElementById("status");
const resultEl = document.getElementById("result");

form.addEventListener("submit", async (e) => {
  e.preventDefault();
  statusEl.textContent = "読み込み中...";
  resultEl.textContent = "";

  try {
    const res = await fetch("/api/menu", { method: "GET" });
    if (!res.ok) throw new Error("HTTP " + res.status);
    const json = await res.json();
    statusEl.textContent = "取得成功";
    resultEl.textContent = JSON.stringify(json, null, 2);
  } catch (err) {
    statusEl.textContent = "エラー: " + err.message;
  }
});
