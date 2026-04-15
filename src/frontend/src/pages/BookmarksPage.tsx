import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import api from '../services/api'
import type { Bookmark, PaginationResponse } from '../types'

function BookmarksPage() {
  const [bookmarks, setBookmarks] = useState<Bookmark[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [page, setPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const pageSize = 10

  useEffect(() => {
    fetchBookmarks(page)
  }, [page])

  const fetchBookmarks = async (pageNum: number) => {
    try {
      setLoading(true)
      setError(null)
      const response = await api.get<PaginationResponse<Bookmark>>('/bookmarks', {
        params: { page: pageNum, pageSize },
      })
      setBookmarks(response.data.data)
      setTotalPages(response.data.totalPages)
    } catch (err) {
      setError('Failed to fetch bookmarks')
      console.error(err)
    } finally {
      setLoading(false)
    }
  }

  const handlePageChange = (newPage: number) => {
    if (newPage >= 1 && newPage <= totalPages) {
      setPage(newPage)
    }
  }

  if (loading) {
    return <div className="bookmarks-page">Loading bookmarks...</div>
  }

  if (error) {
    return <div className="bookmarks-page error">{error}</div>
  }

  return (
    <div className="bookmarks-page">
      <header className="page-header">
        <h1>Bookmarks</h1>
        <Link to="/bookmarks/new" className="btn btn-primary">
          Add Bookmark
        </Link>
      </header>

      <div className="bookmarks-list">
        {bookmarks.length === 0 ? (
          <div className="empty-state">
            <p>No bookmarks yet. Add your first bookmark!</p>
          </div>
        ) : (
          bookmarks.map((bookmark) => (
            <Link key={bookmark.id} to={`/bookmarks/${bookmark.id}`} className="bookmark-card">
              <h3>{bookmark.title}</h3>
              <p className="url">{bookmark.url}</p>
              {bookmark.tags.length > 0 && (
                <div className="tags">
                  {bookmark.tags.map((tag) => (
                    <span key={tag} className="tag-chip">
                      {tag}
                    </span>
                  ))}
                </div>
              )}
            </Link>
          ))
        )}
      </div>

      {totalPages > 1 && (
        <div className="pagination">
          <button
            onClick={() => handlePageChange(page - 1)}
            disabled={page === 1}
            className="btn"
          >
            Previous
          </button>
          <span className="page-indicator">
            Page {page} of {totalPages}
          </span>
          <button
            onClick={() => handlePageChange(page + 1)}
            disabled={page === totalPages}
            className="btn"
          >
            Next
          </button>
        </div>
      )}
    </div>
  )
}

export default BookmarksPage
