(() => {
  document.addEventListener('click', async (event) => {
    const button = event.target.closest('[data-copy-target]');
    if (!button) return;
    const target = document.querySelector(button.dataset.copyTarget);
    if (!target) return;
    try {
      await navigator.clipboard.writeText(target.innerText);
      const previous = button.textContent;
      button.textContent = 'Copied';
      window.setTimeout(() => { button.textContent = previous; }, 1200);
    } catch (_) {
      button.textContent = 'Select manually';
    }
  });

  document.addEventListener('htmx:afterSwap', (event) => {
    if (event.target?.id !== 'component-grid') return;
    document.querySelectorAll('.filter-chip').forEach((chip) => chip.classList.remove('is-active'));
  });
})();
