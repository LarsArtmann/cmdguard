(function () {
  const theme = localStorage.getItem("theme");
  if (theme === "light") {
    document.documentElement.classList.add("light");
  }
})();
