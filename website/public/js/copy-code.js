(function () {
  const btn = document.getElementById("copy-btn");
  if (!btn) return;

  btn.addEventListener("click", function () {
    const code = btn.getAttribute("data-code");
    if (!code) return;

    navigator.clipboard.writeText(code).then(function () {
      const original = btn.textContent;
      btn.textContent = "Copied!";
      btn.classList.add("text-success");
      setTimeout(function () {
        btn.textContent = original;
        btn.classList.remove("text-success");
      }, 2000);
    });
  });
})();
