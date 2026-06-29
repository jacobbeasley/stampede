// Theme Manager - Handles theme selection and persistence using localStorage
(function() {
    'use strict';
    const ThemeManager = {
        storageKey: 'landing-page-theme',
        defaultTheme: 'dark',
        availableThemes: [
            'light', 'dark', 'cupcake', 'bumblebee', 'emerald', 'corporate',
            'synthwave', 'retro', 'cyberpunk', 'valentine', 'halloween',
            'garden', 'forest', 'aqua', 'lofi', 'pastel', 'fantasy',
            'wireframe', 'business', 'acid', 'lemonade', 'night', 'coffee', 'winter',
            'dim', 'nord', 'sunset'
        ],
        init() {
            const savedTheme = localStorage.getItem(this.storageKey);
            const themeToUse = savedTheme && this.availableThemes.includes(savedTheme) ? savedTheme : this.defaultTheme;
            this.setTheme(themeToUse, false);
            this.updateThemeDisplay();
        },
        setTheme(theme, save = true) {
            if (!this.availableThemes.includes(theme)) theme = this.defaultTheme;
            document.documentElement.setAttribute('data-theme', theme);
            if (save) localStorage.setItem(this.storageKey, theme);
            this.updateThemeDisplay();
            window.dispatchEvent(new CustomEvent('themeChanged', { detail: { theme } }));
        },
        updateThemeDisplay() {
            const themeSelects = document.querySelectorAll('.theme-selector');
            themeSelects.forEach(select => {
                const currentTheme = document.documentElement.getAttribute('data-theme');
                select.value = currentTheme || this.defaultTheme;
            });
        }
    };
    window.setTheme = function(theme) { ThemeManager.setTheme(theme); };
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', () => ThemeManager.init());
    } else {
        ThemeManager.init();
    }
})();
