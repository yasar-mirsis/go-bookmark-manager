import { useParams, Link } from 'react-router-dom'
import { useState, useEffect } from 'react'
import api from '../services/api'
import type { Bookmark } from '../types'

function BookmarkDetailPage() {
  const { id } = useParams<{ id: string }>()
  const [bookmark, setBookmark] = useState<Bookmark | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (id) {
      fetchBookmark(id)
    }
  }, [id])

  const fetchBookmark = async (bookmarkId: string) => {
    try {
      setLoading(true)
      setError(null)
      const response = await api.get<Bookmark>(`/bookmarks/${bookmarkId}`)
      setBookmark(response.data)
    } catch (err) {
      setError('Failed to fetch bookmark')
      console.error(err)
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return <div className="bookmark-detail-page">Loading...</div>
  }

  if (error || !bookmark) {
    return (
      <div className="bookmark-detail-page error">
        <p>{error || 'Bookmark not found'}</p>
        <Link to="/" className="btn">
          Back to Bookmarks
        </Link>
      </div>
    )
  }

  return (
    <div className="bookmark-detail-page">
      <Link to="/" className="btn">
        Back to Bookmarks
      </Link>

      <div className="bookmark-detail">
        <h1>{bookmark.title}</h1>
        <a href={bookmark.url} target="_blank" rel="noopener noreferrer" className="url-link">
          {bookmark.url}
        </a>

        {bookmark.description && (
          <div className="description">
            <h2>Description</h2>
            <p>{bookmark.description}</p>
          </div>
        )}

        {bookmark.tags.length > 0 && (
          <div className="tags">
            <h2>Tags</h2>
            <div className="tag-list">
              {bookmark.tags.map((tag) => (
                <span key={tag} className="tag-chip">
                  {tag}
                </span>
              ))}
            </div>
          </div>
        )}

        <div className="metadata">
          <p>
            Created: {new Date(bookmark.createdAt).toLocaleString()}
          </p>
          <p>
            Updated: {new Date(bookmark.updatedAt).toLocaleString()}
          </p>
        </div>
      </div>
    </div>
  )
}

export default BookmarkDetailPage
