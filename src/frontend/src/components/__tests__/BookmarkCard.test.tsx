import { render, screen, fireEvent } from '@testing-library/react'
import '@testing-library/jest-dom'
import BookmarkCard from '../BookmarkCard'
import type { Bookmark } from '../../types'

const mockBookmark: Bookmark = {
  id: '1',
  url: 'https://example.com',
  title: 'Example Bookmark',
  description: 'This is a test bookmark description',
  tags: ['react', 'typescript', 'testing'],
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-01-01T00:00:00Z',
}

const defaultProps = {
  bookmark: mockBookmark,
  onEdit: vi.fn(),
  onDelete: vi.fn(),
  onClick: vi.fn(),
  onTagClick: vi.fn(),
}

describe('BookmarkCard', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders correctly with all props', () => {
    render(<BookmarkCard {...defaultProps} />)

    // Check title is rendered
    expect(screen.getByText('Example Bookmark')).toBeInTheDocument()

    // Check URL is rendered
    expect(screen.getByTitle('https://example.com')).toBeInTheDocument()

    // Check description is rendered
    expect(screen.getByTitle('This is a test bookmark description')).toBeInTheDocument()

    // Check tags are rendered as chips
    expect(screen.getByText('react')).toBeInTheDocument()
    expect(screen.getByText('typescript')).toBeInTheDocument()
    expect(screen.getByText('testing')).toBeInTheDocument()

    // Check action buttons are rendered
    expect(screen.getByRole('button', { name: 'Edit' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Delete' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'View' })).toBeInTheDocument()
  })

  it('truncates long titles correctly', () => {
    const longTitleBookmark: Bookmark = {
      ...mockBookmark,
      title: 'This is a very long title that should be truncated because it exceeds the maximum character limit for display purposes',
    }

    render(<BookmarkCard {...defaultProps} bookmark={longTitleBookmark} />)

    const titleElement = screen.getByTitle(longTitleBookmark.title)
    expect(titleElement.textContent).toContain('...')
  })

  it('truncates long URLs correctly', () => {
    const longUrlBookmark: Bookmark = {
      ...mockBookmark,
      url: 'https://very-long-domain-name-with-many-subdomains.example.com/very/long/path/to/resource',
    }

    render(<BookmarkCard {...defaultProps} bookmark={longUrlBookmark} />)

    const urlElement = screen.getByTitle(longUrlBookmark.url)
    expect(urlElement.textContent).toContain('...')
  })

  it('truncates long descriptions correctly', () => {
    const longDescBookmark: Bookmark = {
      ...mockBookmark,
      description: 'This is a very long description that should be truncated because it exceeds the maximum character limit for display purposes and should show ellipsis at the end',
    }

    render(<BookmarkCard {...defaultProps} bookmark={longDescBookmark} />)

    const descElement = screen.getByTitle(longDescBookmark.description)
    expect(descElement.textContent).toContain('...')
  })

  it('does not render description when empty', () => {
    const noDescBookmark: Bookmark = {
      ...mockBookmark,
      description: '',
    }

    render(<BookmarkCard {...defaultProps} bookmark={noDescBookmark} />)

    // Description should not be in the document
    expect(screen.queryByTitle(noDescBookmark.description)).not.toBeInTheDocument()
  })

  it('does not render tags section when no tags', () => {
    const noTagsBookmark: Bookmark = {
      ...mockBookmark,
      tags: [],
    }

    render(<BookmarkCard {...defaultProps} bookmark={noTagsBookmark} />)

    // No tag chips should be rendered
    expect(screen.queryByText('react')).not.toBeInTheDocument()
  })

  it('fires onClick event when title is clicked', () => {
    render(<BookmarkCard {...defaultProps} />)

    const titleElement = screen.getByText('Example Bookmark')
    fireEvent.click(titleElement)

    expect(defaultProps.onClick).toHaveBeenCalledWith('1')
  })

  it('fires onEdit event when Edit button is clicked', () => {
    render(<BookmarkCard {...defaultProps} />)

    const editButton = screen.getByRole('button', { name: 'Edit' })
    fireEvent.click(editButton)

    expect(defaultProps.onEdit).toHaveBeenCalled()
  })

  it('fires onDelete event when Delete button is clicked and confirmed', () => {
    vi.spyOn(window, 'confirm').mockReturnValue(true)

    render(<BookmarkCard {...defaultProps} />)

    const deleteButton = screen.getByRole('button', { name: 'Delete' })
    fireEvent.click(deleteButton)

    expect(defaultProps.onDelete).toHaveBeenCalledWith('1')

    vi.restoreAllMocks()
  })

  it('does not fire onDelete event when delete is cancelled', () => {
    vi.spyOn(window, 'confirm').mockReturnValue(false)

    render(<BookmarkCard {...defaultProps} />)

    const deleteButton = screen.getByRole('button', { name: 'Delete' })
    fireEvent.click(deleteButton)

    expect(defaultProps.onDelete).not.toHaveBeenCalled()

    vi.restoreAllMocks()
  })

  it('fires onTagClick event when tag chip is clicked', () => {
    render(<BookmarkCard {...defaultProps} />)

    const reactTag = screen.getByText('react')
    fireEvent.click(reactTag)

    expect(defaultProps.onTagClick).toHaveBeenCalledWith('react')
  })

  it('fires onClick event when View button is clicked', () => {
    render(<BookmarkCard {...defaultProps} />)

    const viewButton = screen.getByRole('button', { name: 'View' })
    fireEvent.click(viewButton)

    expect(defaultProps.onClick).toHaveBeenCalledWith('1')
  })

  it('prevents event propagation on Edit button click', () => {
    render(<BookmarkCard {...defaultProps} />)

    const editButton = screen.getByRole('button', { name: 'Edit' })
    fireEvent.click(editButton)

    expect(defaultProps.onClick).not.toHaveBeenCalled()
  })

  it('prevents event propagation on Delete button click', () => {
    vi.spyOn(window, 'confirm').mockReturnValue(false)

    render(<BookmarkCard {...defaultProps} />)

    const deleteButton = screen.getByRole('button', { name: 'Delete' })
    fireEvent.click(deleteButton)

    expect(defaultProps.onClick).not.toHaveBeenCalled()

    vi.restoreAllMocks()
  })

  it('prevents event propagation on tag chip click', () => {
    render(<BookmarkCard {...defaultProps} />)

    const reactTag = screen.getByText('react')
    fireEvent.click(reactTag)

    expect(defaultProps.onClick).not.toHaveBeenCalled()
  })

  // Snapshot test
  it('matches snapshot', () => {
    const { container } = render(<BookmarkCard {...defaultProps} />)
    expect(container).toMatchSnapshot()
  })
})
