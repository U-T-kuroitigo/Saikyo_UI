// 味の名前とAPIが要求するIDを対応付けるための定数
const FLAVOR_TO_API_ID = {
	いちご: "giiku-sai",
	メロン: "giiku-haku",
	ブルーハワイ: "giiku-ten",
	オレンジ: "giiku-camp",
};

// 注文状況を定期的に確認するためのタイマーIDを保持する変数
let statusPollingInterval = null;

/**
 * 注文状況をAPIに問い合わせて、画面表示を更新します。
 * @param {string} orderId - 確認対象の注文IDです。
 * @param {HTMLElement} statusElement - 状況を表示するためのDOM要素です。
 */
async function checkOrderStatus(orderId, statusElement) {
	try {
		const response = await fetch(
			// `https://kakigori-api.fly.dev/v1/stores/UQHVDAEW/orders/${orderId}`
			`/api/orders/${orderId}`
		);
		if (!response.ok) {
			// APIからのエラー応答はコンソールに出力するのみで、画面表示は変更しない
			console.error("Status check failed:", response.status);
			return;
		}
		const result = await response.json();

		// APIから返されたstatusに応じて表示を切り替える
		if (result.status === "pending") {
			statusElement.textContent = "準備中...";
		} else if (result.status === "waitingPickup") {
			statusElement.textContent = "作成完了！";
			// 完了したら定期確認を停止する
			clearInterval(statusPollingInterval);
		}
	} catch (error) {
		console.error("Status check fetch error:", error);
	}
}

/**
 * 注文状況の定期的な確認を開始します。
 * @param {string} orderId - 確認対象の注文IDです。
 * @param {HTMLElement} statusElement - 状況を表示するためのDOM要素です。
 */
function startStatusPolling(orderId, statusElement) {
	// 最初に一度すぐに状況を確認
	checkOrderStatus(orderId, statusElement);
	// その後、10秒ごとに繰り返し確認
	statusPollingInterval = setInterval(
		() => checkOrderStatus(orderId, statusElement),
		10000
	);
}

/**
 * フォームの送信を処理し、APIを呼び出します。
 * @param {Event} event - イベントオブジェクトです。
 * @param {object} elements - 操作するDOM要素の集まりです。
 */
async function handleOrderSubmit(event, elements) {
	event.preventDefault(); // フォームのデフォルトの送信動作をキャンセル

	const formData = new FormData(elements.iceForm);
	const selectedFlavor = formData.get("flavor");

	if (!selectedFlavor) {
		alert("かき氷の味を選んでください。");
		return;
	}

	const menuItemId = FLAVOR_TO_API_ID[selectedFlavor];
	if (!menuItemId) {
		alert("不正な味が選択されました。ページをリロードしてください。");
		return;
	}

	// UIを「注文中」の状態に更新
	elements.orderBtn.disabled = true;
	elements.orderBtn.textContent = "注文中...";

	try {
		// サーバーのAPIエンドポイントに注文データを送信
		const response = await fetch("/api/orders", {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({ menu_item_id: menuItemId }),
		});

		const result = await response.json(); // サーバーからの応答をJSONとして解釈

		if (response.ok) {
			// 成功した場合: フォームを非表示にし、受け取り番号と状況を表示
			elements.orderSection.innerHTML = `
        <h2>ご注文ありがとうございます</h2>
        <p>受け取り番号: <strong>${result.order_number}</strong></p>
        <p>状況: <span id="order-status-text"></span></p>
      `;
			// 新しく作成した状況表示用の要素を取得
			const statusElement = document.getElementById("order-status-text");
			// 注文状況の定期確認を開始
			startStatusPolling(result.id, statusElement);
		} else {
			// 失敗した場合
			elements.resultDiv.textContent = `エラー: ${
				result.error || "注文に失敗しました。"
			}`;
			elements.resultDiv.style.color = "red";
			elements.orderBtn.disabled = false;
			elements.orderBtn.textContent = "決定";
		}
	} catch (error) {
		// 通信自体に失敗した場合
		elements.resultDiv.textContent = "通信エラーが発生しました。";
		elements.resultDiv.style.color = "red";
		elements.orderBtn.disabled = false;
		elements.orderBtn.textContent = "決定";
		console.error("Fetch Error:", error);
	}
}

/**
 * かき氷注文フォームの初期化処理を行います。
 */
function initializeIceOrderForm() {
	// HTMLから必要な要素を取得
	const elements = {
		iceForm: document.getElementById("ice-form"),
		orderBtn: document.getElementById("order-btn"),
		resultDiv: document.getElementById("order-result"),
		orderSection: document.getElementById("ice-order-section"),
	};

	// フォームが存在する場合のみ、ボタンにクリックイベントを設定
	if (elements.iceForm) {
		elements.orderBtn.addEventListener("click", (event) =>
			handleOrderSubmit(event, elements)
		);
	}
}

// ページのHTMLがすべて読み込まれたら、初期化処理を実行
document.addEventListener("DOMContentLoaded", initializeIceOrderForm);
