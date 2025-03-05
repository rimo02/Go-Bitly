import {RouterProvider, createBrowserRouter } from 'react-router-dom'
import Home from './pages/home';
import ShortenPage from './pages/shortenPage';
const router = createBrowserRouter([
  {
    path: "/",
    children: [
      { path: "/", element: <ShortenPage /> },
      { path: "/home", element: <Home /> }
    ]
  }
])

function App() {
  return (
    <RouterProvider router={router} />
  )
}

export default App;
