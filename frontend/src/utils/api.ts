// API configuration utility
export const getApiBaseUrl = (): string => {
  // In production (Docker), frontend and backend are served from the same origin
  // so we can use relative paths
  if (window.location.hostname !== 'localhost' || process.env.NODE_ENV === 'production') {
    return '/api';
  }
  // In development, backend runs on port 8080
  return 'http://localhost:8080/api';
};

export const API_BASE = getApiBaseUrl();