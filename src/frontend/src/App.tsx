import { Routes, Route } from 'react-router-dom'
import BookmarkPage from './pages/BookmarkPage'
import BookmarkDetailPage from './pages/BookmarkDetailPage'

function App() {
  return (
    <div className="app">
      <Routes>
        <Route path="/" element={<BookmarkPage />} />
        <Route path="/bookmarks/:id" element={<BookmarkDetailPage />} />
      </Routes>
    </div>
  )
}

export default App
