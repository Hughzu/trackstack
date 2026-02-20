type SignedFetchInit = RequestInit & { body?: BodyInit | Record<string, unknown> | null };

declare global {
  interface Window {
    signedFetch?: (input: RequestInfo | URL, init?: SignedFetchInit) => Promise<Response>;
  }
}

const ensureSignedFetch = () => {
  if (window.signedFetch) return;

  window.signedFetch = async (input: RequestInfo | URL, init: SignedFetchInit = {}) => {
    const method = (init.method || "GET").toUpperCase();
    if (method === "GET" || method === "HEAD") {
      return fetch(input, init);
    }

    const headers = new Headers(init.headers || {});
    let body = init.body ?? null;
    let bodyString = "";

    if (typeof body === "string") {
      bodyString = body;
    } else if (body instanceof URLSearchParams) {
      bodyString = body.toString();
      if (!headers.has("Content-Type")) {
        headers.set("Content-Type", "application/x-www-form-urlencoded");
      }
      body = bodyString;
    } else if (body instanceof FormData) {
      bodyString = new URLSearchParams(body).toString();
      if (!headers.has("Content-Type")) {
        headers.set("Content-Type", "application/x-www-form-urlencoded");
      }
      body = bodyString;
    } else if (body && typeof body === "object") {
      bodyString = JSON.stringify(body);
      if (!headers.has("Content-Type")) {
        headers.set("Content-Type", "application/json");
      }
      body = bodyString;
    }

    if (window.crypto?.subtle) {
      const digest = await window.crypto.subtle.digest(
        "SHA-256",
        new TextEncoder().encode(bodyString)
      );
      const hash = Array.from(new Uint8Array(digest))
        .map((b) => b.toString(16).padStart(2, "0"))
        .join("");
      headers.set("x-amz-content-sha256", hash);
    }

    return fetch(input, { ...init, headers, body });
  };
};

const initApiForms = () => {
  const forms = Array.from(document.querySelectorAll<HTMLFormElement>("form[data-api-form]"));

  forms.forEach((form) => {
    if (form.dataset.apiFormBound === "true") return;
    form.dataset.apiFormBound = "true";

    form.addEventListener("submit", async (event) => {
      event.preventDefault();

      const formData = new FormData(form);
      const payload = Object.fromEntries(formData);
      const action = form.getAttribute("action") || window.location.href;
      const method = form.getAttribute("method")?.toUpperCase() || "POST";
      const redirectUrl = form.dataset.redirect;

      try {
        const fetcher = window.signedFetch || window.fetch;
        const response = await fetcher(action, {
          method,
          headers: {
            "Content-Type": "application/json"
          },
          body: JSON.stringify(payload)
        });

        if (response.ok) {
          if (redirectUrl) window.location.href = redirectUrl;
          else window.location.reload();
          return;
        }

        const errorTargetId = form.dataset.errorTarget;
        if (errorTargetId) {
          const errEl = document.getElementById(errorTargetId);
          if (errEl) {
            let msg = "Request failed.";
            try {
              const data = await response.json();
              if (data?.error) msg = data.error;
            } catch {
              // noop
            }
            errEl.textContent = msg;
            errEl.classList.remove("hidden");
            return;
          }
        }
      } catch (error) {
        console.error("SigV4 signing or network error:", error);
      }

      const errorTargetId = form.dataset.errorTarget;
      if (errorTargetId) {
        const errEl = document.getElementById(errorTargetId);
        if (errEl) {
          errEl.textContent = "Network error occurred.";
          errEl.classList.remove("hidden");
          return;
        }
      }

      window.alert("Failed to process request.");
    });
  });
};

const dropdownState = {
  roots: new Set<HTMLElement>(),
  globalBound: false
};

const getDropdownPanel = (root: HTMLElement) => root.querySelector<HTMLElement>("[data-menu-panel]");
const getDropdownTrigger = (root: HTMLElement) => root.querySelector<HTMLElement>("[data-menu-trigger]");
const getDropdownControl = (root: HTMLElement) => getDropdownTrigger(root)?.firstElementChild as HTMLElement | null;

const isDropdownOpen = (root: HTMLElement) => {
  const panel = getDropdownPanel(root);
  return panel && !panel.classList.contains("hidden");
};

const setDropdownOpen = (root: HTMLElement, open: boolean) => {
  const panel = getDropdownPanel(root);
  if (!panel) return;
  panel.classList.toggle("hidden", !open);
  const control = getDropdownControl(root);
  if (control) {
    control.setAttribute("aria-expanded", String(open));
  }
};

const closeAllDropdowns = (except?: HTMLElement) => {
  dropdownState.roots.forEach((root) => {
    if (root !== except) {
      setDropdownOpen(root, false);
    }
  });
};

const initDropdownMenus = () => {
  dropdownState.roots.forEach((root) => {
    if (!root.isConnected) dropdownState.roots.delete(root);
  });

  const roots = Array.from(document.querySelectorAll<HTMLElement>("[data-menu-root]"));

  roots.forEach((root) => {
    if (root.dataset.menuBound === "true") return;
    root.dataset.menuBound = "true";
    dropdownState.roots.add(root);

    const trigger = getDropdownTrigger(root);
    trigger?.addEventListener("click", (event) => {
      event.stopPropagation();
      const next = !isDropdownOpen(root);
      closeAllDropdowns(root);
      setDropdownOpen(root, next);
    });
  });

  if (!dropdownState.globalBound) {
    dropdownState.globalBound = true;

    document.addEventListener("click", (event) => {
      dropdownState.roots.forEach((root) => {
        if (!isDropdownOpen(root)) return;
        if (!root.contains(event.target as Node)) {
          setDropdownOpen(root, false);
        }
      });
    });

    document.addEventListener("keydown", (event) => {
      if (event.key !== "Escape") return;
      dropdownState.roots.forEach((root) => setDropdownOpen(root, false));
    });
  }
};

type ConfirmModalConfig = {
  dialog: HTMLDialogElement;
  triggerAttribute: string;
  idAttribute?: string;
  endpoint: string;
  method: string;
  payloadKey: string;
  sendPayload: boolean;
  errorMessage: string;
  successRedirect?: string;
  confirmLoadingLabel: string;
  confirmBtn: HTMLButtonElement | null;
  closeBtns: NodeListOf<HTMLElement>;
  pendingId: string | null;
};

const confirmModalState = {
  modals: [] as ConfirmModalConfig[],
  globalBound: false
};

const initConfirmModals = () => {
  confirmModalState.modals = confirmModalState.modals.filter((modal) => modal.dialog.isConnected);

  const dialogs = Array.from(document.querySelectorAll<HTMLDialogElement>("dialog[data-confirm-dialog]"));

  dialogs.forEach((dialog) => {
    if (dialog.dataset.confirmBound === "true") return;
    dialog.dataset.confirmBound = "true";

    const triggerAttribute = dialog.dataset.triggerAttribute || "";
    if (!triggerAttribute) return;

    const config: ConfirmModalConfig = {
      dialog,
      triggerAttribute,
      idAttribute: dialog.dataset.idAttribute,
      endpoint: dialog.dataset.endpoint || "",
      method: (dialog.dataset.method || "DELETE").toUpperCase(),
      payloadKey: dialog.dataset.payloadKey || "id",
      sendPayload: dialog.dataset.sendPayload !== "false",
      errorMessage: dialog.dataset.errorMessage || "Something went wrong. Please try again.",
      successRedirect: dialog.dataset.successRedirect || undefined,
      confirmLoadingLabel: dialog.dataset.confirmLoadingLabel || "Working...",
      confirmBtn: dialog.querySelector<HTMLButtonElement>("[data-confirm-modal]"),
      closeBtns: dialog.querySelectorAll<HTMLElement>("[data-close-modal]"),
      pendingId: null
    };

    config.closeBtns.forEach((btn) => {
      btn.addEventListener("click", () => dialog.close());
    });

    config.confirmBtn?.addEventListener("click", async () => {
      if (!config.pendingId) return;

      const originalText = config.confirmBtn?.textContent || "";
      if (config.confirmBtn) {
        config.confirmBtn.disabled = true;
        config.confirmBtn.textContent = config.confirmLoadingLabel;
      }

      try {
        let resolvedEndpoint = config.endpoint;
        const requestInit: RequestInit = {
          method: config.method,
          headers: { "Content-Type": "application/json" }
        };

        if (config.sendPayload) {
          if (config.method === "DELETE") {
            const url = new URL(config.endpoint, window.location.origin);
            url.searchParams.set(config.payloadKey, config.pendingId);
            resolvedEndpoint = url.toString();
          } else {
            requestInit.body = JSON.stringify({ [config.payloadKey]: config.pendingId });
          }
        }

        const signedFetch = window.signedFetch || window.fetch;
        const response = await signedFetch(resolvedEndpoint, requestInit);

        if (response.ok) {
          if (config.successRedirect) {
            window.location.href = config.successRedirect;
            return;
          }
          window.location.reload();
          return;
        }

        throw new Error("Request failed");
      } catch (err) {
        console.error(err);
        window.alert(config.errorMessage);
      } finally {
        if (config.confirmBtn) {
          config.confirmBtn.disabled = false;
          config.confirmBtn.textContent = originalText;
        }
        dialog.close();
      }
    });

    confirmModalState.modals.push(config);
  });

  if (!confirmModalState.globalBound) {
    confirmModalState.globalBound = true;

    document.addEventListener("click", (event) => {
      const target = event.target as Element | null;
      if (!target) return;

      for (const modal of confirmModalState.modals) {
        if (!modal.dialog.isConnected) continue;
        const trigger = target.closest(`[${modal.triggerAttribute}]`);
        if (!trigger) continue;

        if (modal.sendPayload) {
          if (!modal.idAttribute) return;
          modal.pendingId = trigger.getAttribute(modal.idAttribute);
          if (!modal.pendingId) return;
        } else {
          modal.pendingId = "__no_payload__";
        }

        modal.dialog.showModal();
        return;
      }
    });
  }
};

const initDialogBackdropClose = () => {
  const dialogs = Array.from(document.querySelectorAll<HTMLDialogElement>("dialog[data-modal-root]"));

  dialogs.forEach((dialog) => {
    if (dialog.dataset.modalBound === "true") return;
    dialog.dataset.modalBound = "true";

    dialog.addEventListener("click", (event) => {
      if (event.target === dialog) {
        dialog.close();
      }
    });
  });
};

const initUi = () => {
  ensureSignedFetch();
  initApiForms();
  initDropdownMenus();
  initConfirmModals();
  initDialogBackdropClose();
};

const bindUi = () => {
  initUi();
};

if (document.readyState === "loading") {
  document.addEventListener("DOMContentLoaded", bindUi);
} else {
  bindUi();
}

document.addEventListener("astro:page-load", bindUi);
document.addEventListener("astro:after-swap", bindUi);

export {};
