// 味の名前とAPIが要求するIDを対応付けるための定数
const FLAVOR_TO_API_ID = {
	いちご: "giiku-sai",
	メロン: "giiku-haku",
	ブルーハワイ: "giiku-ten",
	オレンジ: "giiku-camp",
};

/**
 * フォームのUI（見た目）の状態を更新します。
 * @param {object} elements - 操作するDOM要素の集まりです。
 * @param {boolean} isLoading - ローディング状態（true）か通常状態（false）かを指定します。
 * @param {string} [message] - 表示するメッセージです。
 * @param {boolean} [isError] - メッセージがエラーかどうかを示します。
 */
function updateUIState(elements, isLoading, message = "", isError = false) {
	elements.orderBtn.disabled = isLoading;
	elements.orderBtn.textContent = isLoading ? "注文中..." : "決定";

	if (message) {
		elements.resultDiv.textContent = message;
		elements.resultDiv.style.color = isError ? "red" : "black";
	} else {
		elements.resultDiv.textContent = "";
	}
}

/**
 * 注文が成功したときの画面を表示します。
 * @param {object} elements - 操作するDOM要素の集まりです。
 * @param {string} message - 表示する成功メッセージです。
 */
function showSuccessUI(elements, message) {
	elements.orderSection.innerHTML = `<h2>注文完了！</h2><p>${message}</p>`;
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

	updateUIState(elements, true); // UIを「注文中」の状態に更新

	try {
		// サーバーのAPIエンドポイントに注文データを送信
		const response = await fetch("/api/orders", {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({ menu_item_id: menuItemId }),
		});

		const result = await response.json(); // サーバーからの応答をJSONとして解釈

		if (response.ok) {
			// 成功した場合
			showSuccessUI(elements, result.message || "注文を受け付けました。");
		} else {
			// 失敗した場合
			updateUIState(
				elements,
				false,
				`エラー: ${result.error || "注文に失敗しました。"}`,
				true
			);
		}
	} catch (error) {
		// 通信自体に失敗した場合
		updateUIState(elements, false, "通信エラーが発生しました。", true);
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
