export const formatIndonesianTimestamp = (isoTimestamp) => {
    const date = new Date(isoTimestamp);
  
    // Check if the date is valid
    if (isNaN(date.getTime())) {
        return "Invalid date";
    }
  
    // Format the date to Indonesian locale
    const options = {
        weekday: 'long',      // "Senin", "Selasa", etc.
        year: 'numeric',      // Full year
        month: 'long',        // Full month name
        day: 'numeric',       // Day of the month
        hour: '2-digit',      // Hours with leading zero
        timeZone: 'Asia/Jakarta',
        timeZoneName: 'short' ,
        minute: '2-digit',     // Minutes with leading zero
        second: '2-digit',
    };
  
    return new Intl.DateTimeFormat('id-ID', options).format(date);
  }