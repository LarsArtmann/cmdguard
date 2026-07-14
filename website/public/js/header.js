(function () {
  const toggle = document.getElementById("theme-toggle");
  if (!toggle) return;

  const lightIcon = toggle.querySelector(".theme-icon-light");
  const darkIcon = toggle.querySelector(".theme-icon-dark");

  function updateIcons() {
    if (document.documentElement.classList.contains("light")) {
      if (lightIcon) lightIcon.classList.add("hidden");
      if (darkIcon) darkIcon.classList.remove("hidden");
    } else {
      if (lightIcon) lightIcon.classList.remove("hidden");
      if (darkIcon) darkIcon.classList.add("hidden");
    }
  }

  updateIcons();

  toggle.addEventListener("click", function () {
    document.documentElement.classList.toggle("light");
    const isLight = document.documentElement.classList.contains("light");
    localStorage.setItem("theme", isLight ? "light" : "dark");
    updateIcons();
  });
})();
