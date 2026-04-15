import { Routes, Route } from 'react-router-dom'
import BookmarksPage from './pages/BookmarksPage'
import BookmarkDetailPage from './pages/BookmarkDetailPage'

function App() {
  return (
    <div className="app">
      <Routes>
        <Route path="/" element={<BookmarksPage />} />
        <Route path="/bookmarks/:id" element={<BookmarkDetailPage />} />
      </Routes>
    </div>
  )
}

export default App
