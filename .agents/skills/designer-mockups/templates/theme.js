// Theme Manager - Handles theme selection and persistence using sessionStorage
// Include this script in all pages: <script src="theme.js"></script>

(function() {
    'use strict';

    // Theme manager
    const ThemeManager = {
        storageKey: 'todoflow-theme',
        defaultTheme: 'dark', // Set default fallback theme to dark to match the application guidelines
        availableThemes: [
            'light', 'dark', 'cupcake', 'bumblebee', 'emerald', 'corporate',
            'synthwave', 'retro', 'cyberpunk', 'valentine', 'halloween',
            'garden', 'forest', 'aqua', 'lofi', 'pastel', 'fantasy',
            'wireframe', 'business', 'acid', 'lemonade', 'night', 'coffee', 'winter'
        ],

        init() {
            // Load saved theme or use default
            const savedTheme = sessionStorage.getItem(this.storageKey);
            const themeToUse = savedTheme && this.availableThemes.includes(savedTheme) 
                ? savedTheme 
                : this.defaultTheme;
            
            this.setTheme(themeToUse, false);
            this.updateThemeDisplay();
        },

        setTheme(theme, save = true) {
            if (!this.availableThemes.includes(theme)) {
                console.warn(`Theme "${theme}" is not available. Using default.`);
                theme = this.defaultTheme;
            }

            // Set the theme on the html element
            document.documentElement.setAttribute('data-theme', theme);

            // Save to sessionStorage
            if (save) {
                sessionStorage.setItem(this.storageKey, theme);
            }

            // Update theme display in UI
            this.updateThemeDisplay();

            // Dispatch custom event for other scripts to listen to
            window.dispatchEvent(new CustomEvent('themeChanged', { detail: { theme } }));
        },

        updateThemeDisplay() {
            // Update any theme display elements
            const themeDisplays = document.querySelectorAll('#current-theme');
            themeDisplays.forEach(display => {
                const currentTheme = document.documentElement.getAttribute('data-theme');
                display.textContent = currentTheme || this.defaultTheme;
            });

            // Update any select elements
            const themeSelects = document.querySelectorAll('select[onchange*="setTheme"]');
            themeSelects.forEach(select => {
                const currentTheme = document.documentElement.getAttribute('data-theme');
                select.value = currentTheme || this.defaultTheme;
            });
        }
    };

    // Expose setTheme function globally for onclick handlers
    window.setTheme = function(theme) {
        ThemeManager.setTheme(theme);
    };

    // Initialize on DOM ready
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', () => ThemeManager.init());
    } else {
        ThemeManager.init();
    }

    // Log current theme for debugging
    console.log(`TodoFlow Theme Manager initialized. Current theme: ${document.documentElement.getAttribute('data-theme')}`);
})();
