// Newsletter Widget JavaScript
function showNewsletterPopup() {
    document.getElementById('newsletterOverlay').classList.add('active');
    document.body.style.overflow = 'hidden';
}

function hideNewsletterPopup() {
    document.getElementById('newsletterOverlay').classList.remove('active');
    document.body.style.overflow = 'auto';
}

// Initialize event listeners when DOM is ready
function initializeNewsletterEventListeners() {
    // Close popup when clicking outside
    const overlay = document.getElementById('newsletterOverlay');
    if (overlay) {
        overlay.addEventListener('click', function(e) {
            if (e.target === this) {
                hideNewsletterPopup();
            }
        });
    }

    // Close popup with Escape key
    document.addEventListener('keydown', function(e) {
        if (e.key === 'Escape') {
            hideNewsletterPopup();
        }
    });

    // Handle form submission
    const form = document.getElementById('newsletterForm');
    if (form) {
        form.addEventListener('submit', async function(e) {
            e.preventDefault();
            
            const form = e.target;
            const submitButton = form.querySelector('.newsletter-submit');
            const messageDiv = document.getElementById('newsletterMessage');
            const formData = new FormData(form);

            // Disable submit button and show loading state
            submitButton.disabled = true;
            submitButton.textContent = 'Subscribing...';
            messageDiv.style.display = 'none';

            try {
                const response = await fetch('https://api.web3forms.com/submit', {
                    method: 'POST',
                    body: formData
                });

                const result = await response.json();

                if (response.ok && result.success) {
                    messageDiv.className = 'newsletter-message success';
                    messageDiv.textContent = 'Thank you for subscribing! Check your email for confirmation.';
                    messageDiv.style.display = 'block';
                    form.reset();
                    // Re-check the newsletter checkbox
                    document.getElementById('newsletterSignup').checked = true;
                    
                    // Auto-close popup after 3 seconds
                    setTimeout(() => {
                        hideNewsletterPopup();
                    }, 3000);
                } else {
                    throw new Error(result.message || 'Subscription failed');
                }
            } catch (error) {
                messageDiv.className = 'newsletter-message error';
                messageDiv.textContent = 'Something went wrong. Please try again.';
                messageDiv.style.display = 'block';
                console.error('Newsletter subscription error:', error);
            } finally {
                // Re-enable submit button
                submitButton.disabled = false;
                submitButton.textContent = 'Subscribe';
            }
        });
    }
}

// Example: Auto-show popup after page loads (customize as needed)
// You can remove this or modify the trigger logic based on your needs
document.addEventListener('DOMContentLoaded', function() {
    // Initialize event listeners first
    initializeNewsletterEventListeners();
    
    // Show popup after 5 seconds
    setTimeout(() => {
        showNewsletterPopup();
    }, 5000);
    
    // Or trigger based on scroll, exit intent, etc.
    // Example: Show on scroll down 50%
    /*
    let hasShown = false;
    window.addEventListener('scroll', function() {
        if (!hasShown && window.scrollY > document.body.scrollHeight * 0.5) {
            hasShown = true;
            showNewsletterPopup();
        }
    });
    */
});
