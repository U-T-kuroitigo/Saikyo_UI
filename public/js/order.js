const select = document.getElementById("menuSelect");
const desc = document.getElementById("menuDesc");
const form = document.getElementById("orderForm");
const statusEl = document.getElementById("status");
const resultEl = document.getElementById("result");
const submitBtn = document.getElementById("submitBtn");

let menuItems = [];

async function loadMenu() {
	select.innerHTML = '<option value="">読み込み中...</option>';
	desc.textContent = "";
	try {
		const res = await fetch("/api/menu");
		if (!res.ok) throw new Error("HTTP " + res.status);
		const data = await res.json();
		menuItems = Array.isArray(data.menu) ? data.menu : [];
		if (menuItems.length === 0) {
			select.innerHTML =
				'<option value="">メニューが見つかりません</option>';
			return;
		}
		select.innerHTML =
			'<option value="">選択してください</option>' +
			menuItems
				.map((m) => `<option value="${m.id}">${m.name}</option>`)
				.join("");
	} catch (e) {
		select.innerHTML = '<option value="">取得に失敗しました</option>';
		statusEl.textContent = "メニュー取得エラー: " + e.message;
	}
}

select.addEventListener("change", () => {
	const id = select.value;
	const found = menuItems.find((m) => m.id === id);
	desc.textContent = found ? found.description || "" : "";
});

form.addEventListener("submit", async (e) => {
	e.preventDefault();
	const id = select.value;
	if (!id) {
		statusEl.textContent = "メニューを選択してください";
		return;
	}

	statusEl.textContent = "送信中...";
	resultEl.textContent = "";
	submitBtn.disabled = true;

	try {
		const res = await fetch("/api/orders", {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({ menu_item_id: id }),
		});
		const text = await res.text();
		try {
			const json = JSON.parse(text);
			resultEl.textContent = JSON.stringify(json, null, 2);
		} catch (_) {
			resultEl.textContent = text;
		}
		statusEl.textContent = res.ok ? "注文作成に成功" : "注文作成でエラー";
	} catch (e) {
		statusEl.textContent = "送信エラー: " + e.message;
	} finally {
		submitBtn.disabled = false;
	}
});

loadMenu();
