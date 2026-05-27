import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import AppRoutes from './routes/index.tsx'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { Toaster } from 'react-hot-toast'
import "leaflet/dist/leaflet.css";

const queryClient = new QueryClient();

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <Toaster
      position="top-right"
      reverseOrder={false}
      toastOptions={{
        duration: 5000,
      }}
      />
      <AppRoutes/>
    </QueryClientProvider>
  </StrictMode>,
)
