function showNotification(message) {
    const notification = document.createElement('div');
    notification.className = 'notification';
    notification.textContent = message;
    document.body.appendChild(notification);

    // Trigger animation
    setTimeout(() => {
        notification.classList.add('show');
    }, 10);

    // Remove after 3 seconds
    setTimeout(() => {
        notification.classList.remove('show');
        setTimeout(() => {
            notification.remove();
        }, 300);
    }, 3000);
}

function copyToClipboard(elementId) {
    const secretValue = document.getElementById(elementId).value;
    navigator.clipboard.writeText(secretValue).then(() => {
        showNotification('✓ Copied to clipboard');
    }).catch(() => {
        showNotification('✗ Failed to copy');
    });
}

function showResult(resultId) {
    const resultElement = document.getElementById(resultId);
    resultElement.hidden = false;
    resultElement.classList.add('fadeIn');
}