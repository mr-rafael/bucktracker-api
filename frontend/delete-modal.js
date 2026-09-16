(function () {
  let dialog;

  function ensureDialog() {
    if (dialog) {
      return dialog;
    }

    dialog = document.createElement("dialog");
    dialog.className = "modal";
    dialog.innerHTML = `
      <div class="modal-box">
        <h3 class="text-lg font-bold" data-modal-title></h3>
        <p class="py-4 whitespace-pre-wrap" data-modal-message></p>
        <div class="modal-action" data-modal-actions></div>
      </div>
    `;
    document.body.appendChild(dialog);
    return dialog;
  }

  function setContent(title, message, actions) {
    const modal = ensureDialog();
    modal.querySelector("[data-modal-title]").textContent = title;
    modal.querySelector("[data-modal-message]").textContent = message;

    const actionsContainer = modal.querySelector("[data-modal-actions]");
    actionsContainer.replaceChildren();
    actions.forEach((action) => actionsContainer.appendChild(action));
  }

  function createButton(label, className, onClick) {
    const button = document.createElement("button");
    button.type = "button";
    button.className = className;
    button.textContent = label;
    button.addEventListener("click", onClick);
    return button;
  }

  function closeModal() {
    ensureDialog().close();
  }

  function showResult(title, message) {
    setContent(title, message, [
      createButton("OK", "btn btn-primary", closeModal),
    ]);
    ensureDialog().showModal();
  }

  function confirmAndDelete({ name, kind, endpoint, onSuccess }) {
    const displayName = name || `unnamed ${kind}`;
    const deleteButton = createButton("Delete", "btn btn-error", async () => {
      deleteButton.disabled = true;
      cancelButton.disabled = true;
      deleteButton.textContent = "Deleting...";

      try {
        const response = await authenticatedApi.fetch(endpoint, {
          method: "DELETE",
        });

        if (!response.ok) {
          showResult("Delete failed", await apiErrors.formatResponse(response));
          return;
        }

        if (typeof onSuccess === "function") {
          onSuccess();
        }

        showResult(
          "Deleted",
          `${kind} \'${displayName}\' Successfully Deleted`,
        );
      } catch (error) {
        if (authenticatedApi.isRedirectError(error)) {
          return;
        }

        const message =
          error.isBackendResponse || error.isInvalidResponse
            ? error.message
            : apiErrors.formatNetworkError(error);
        showResult("Delete failed", message);
      }
    });
    const cancelButton = createButton("Cancel", "btn btn-ghost", closeModal);

    setContent(
      "Confirm delete",
      `Are you sure you want to delete \'${displayName}\'?`,
      [cancelButton, deleteButton],
    );
    ensureDialog().showModal();
  }

  function createTrashButton({ name, kind, endpoint, onSuccess }) {
    const button = document.createElement("button");
    button.type = "button";
    button.className =
      "text-base-content/50 transition-colors hover:text-error";
    button.setAttribute("aria-label", `Delete ${name || kind}`);
    button.title = "Delete";
    button.innerHTML = `
      <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        stroke-width="1.8"
        stroke="currentColor"
        class="size-6"
        aria-hidden="true"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673A2.25 2.25 0 0 1 15.916 21.75H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0"
        />
      </svg>
    `;
    button.addEventListener("click", (event) => {
      event.preventDefault();
      event.stopPropagation();
      confirmAndDelete({ name, kind, endpoint, onSuccess });
    });
    return button;
  }

  window.deleteModal = Object.freeze({
    confirmAndDelete,
    createTrashButton,
  });
})();
