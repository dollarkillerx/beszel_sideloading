// API configuration utility
export const getApiBaseUrl = (): string => {
  // In development, when frontend runs on port 3000, backend runs on port 8080
  if (window.location.port === '3000') {
    return 'http://localhost:8080/api';
  }
  // In production (Docker), frontend and backend are served from the same origin
  return '/api';
};

export const API_BASE = getApiBaseUrl();