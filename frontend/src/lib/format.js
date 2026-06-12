export function formatDate(dateString) {
    if (!dateString) return '-';
    try {
        const date = new Date(dateString);
        return date.toLocaleString('id-ID', {
            year: 'numeric',
            month: 'long',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
            second: '2-digit',
        });
    } catch (e) {
        return dateString;
    }
}

export function formatQuota(quota) {
    return (quota === 0 || quota === -1) ? 'Tanpa Batas' : `${quota} Transaksi`;
}
