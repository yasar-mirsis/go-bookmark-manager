import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import '@testing-library/jest-dom'
import SearchBar from '../SearchBar'

const defaultProps = {
  onSearch: vi.fn(),
  onClear: vi.fn(),
}

describe('SearchBar', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('renders search input correctly', () => {
    render(<SearchBar {...defaultProps} />)

    const input = screen.getByRole('textbox', { name: /Search bookmarks/i })
    expect(input).toBeInTheDocument()
    expect(input).toHaveAttribute('placeholder', 'Search bookmarks by title, description, or URL...')
  })

  it('renders clear button when query has value', () => {
    const { rerender } = render(<SearchBar {...defaultProps} />)

    // Initially no clear button
    expect(screen.queryByRole('button', { name: /Clear search/i })).not.toBeInTheDocument()

    // Type in the search box
    const input = screen.getByRole('textbox', { name: /Search bookmarks/i })
    fireEvent.change(input, { target: { value: 'test' } })

    // Clear button should appear
    expect(screen.getByRole('button', { name: /Clear search/i })).toBeInTheDocument()
  })

  it('does not render clear button when query is empty', () => {
    render(<SearchBar {...defaultProps} />)

    expect(screen.queryByRole('button', { name: /Clear search/i })).not.toBeInTheDocument()
  })

  it('calls onSearch after debounce timer when typing', async () => {
    render(<SearchBar {...defaultProps} />)

    const input = screen.getByRole('textbox', { name: /Search bookmarks/i })
    fireEvent.change(input, { target: { value: 'test' } })

    // onSearch should not be called immediately
    expect(defaultProps.onSearch).not.toHaveBeenCalled()

    // Advance timer by 300ms (debounce time)
    vi.advanceTimersByTime(300)

    // onSearch should be called with the query value
    await waitFor(() => {
      expect(defaultProps.onSearch).toHaveBeenCalledWith('test')
    })
  })

  it('debounces multiple rapid inputs correctly', async () => {
    render(<SearchBar {...defaultProps} />)

    const input = screen.getByRole('textbox', { name: /Search bookmarks/i })

    // Simulate rapid typing
    fireEvent.change(input, { target: { value: 't' } })
    vi.advanceTimersByTime(100)

    fireEvent.change(input, { target: { value: 'te' } })
    vi.advanceTimersByTime(100)

    fireEvent.change(input, { target: { value: 'tes' } })
    vi.advanceTimersByTime(100)

    fireEvent.change(input, { target: { value: 'test' } })
    vi.advanceTimersByTime(300)

    // Should only call onSearch once with the final value
    await waitFor(() => {
      expect(defaultProps.onSearch).toHaveBeenCalledTimes(1)
      expect(defaultProps.onSearch).toHaveBeenCalledWith('test')
    })
  })

  it('clears debounce timer on unmount', () => {
    const { unmount } = render(<SearchBar {...defaultProps} />)

    const input = screen.getByRole('textbox', { name: /Search bookmarks/i })
    fireEvent.change(input, { target: { value: 'test' } })

    // Don't advance timers, just unmount
    unmount()

    // Advance timers after unmount - should not call onSearch
    vi.advanceTimersByTime(300)

    expect(defaultProps.onSearch).not.toHaveBeenCalled()
  })

  it('calls onClear when clear button is clicked', () => {
    render(<SearchBar {...defaultProps} />)

    const input = screen.getByRole('textbox', { name: /Search bookmarks/i })
    fireEvent.change(input, { target: { value: 'test' } })

    const clearButton = screen.getByRole('button', { name: /Clear search/i })
    fireEvent.click(clearButton)

    expect(defaultProps.onClear).toHaveBeenCalled()
    expect(input).toHaveValue('')
  })

  it('clears search when Escape key is pressed', () => {
    render(<SearchBar {...defaultProps} />)

    const input = screen.getByRole('textbox', { name: /Search bookmarks/i })
    fireEvent.change(input, { target: { value: 'test' } })

    fireEvent.keyDown(input, { key: 'Escape', code: 'Escape' })

    expect(defaultProps.onClear).toHaveBeenCalled()
    expect(input).toHaveValue('')
  })

  it('does not call onClear when other keys are pressed', () => {
    render(<SearchBar {...defaultProps} />)

    const input = screen.getByRole('textbox', { name: /Search bookmarks/i })
    fireEvent.change(input, { target: { value: 'test' } })

    fireEvent.keyDown(input, { key: 'Enter', code: 'Enter' })

    expect(defaultProps.onClear).not.toHaveBeenCalled()
  })

  it('updates input value when typing', () => {
    render(<SearchBar {...defaultProps} />)

    const input = screen.getByRole('textbox', { name: /Search bookmarks/i })

    fireEvent.change(input, { target: { value: 'hello world' } })

    expect(input).toHaveValue('hello world')
  })

  it('calls onSearch with empty string when clearing', async () => {
    render(<SearchBar {...defaultProps} />)

    const input = screen.getByRole('textbox', { name: /Search bookmarks/i })
    fireEvent.change(input, { target: { value: 'test' } })
    vi.advanceTimersByTime(300)

    const clearButton = screen.getByRole('button', { name: /Clear search/i })
    fireEvent.click(clearButton)

    // onClear is called, not onSearch with empty string
    expect(defaultProps.onClear).toHaveBeenCalled()
  })

  it('has correct ARIA attributes for accessibility', () => {
    render(<SearchBar {...defaultProps} />)

    const input = screen.getByRole('textbox', { name: /Search bookmarks/i })
    expect(input).toHaveAttribute('aria-label', 'Search bookmarks')
  })

  it('clear button has correct ARIA label', () => {
    render(<SearchBar {...defaultProps} />)

    const input = screen.getByRole('textbox', { name: /Search bookmarks/i })
    fireEvent.change(input, { target: { value: 'test' } })

    const clearButton = screen.getByRole('button', { name: /Clear search/i })
    expect(clearButton).toHaveAttribute('aria-label', 'Clear search')
  })

  it('handles special characters in search query', async () => {
    render(<SearchBar {...defaultProps} />)

    const input = screen.getByRole('textbox', { name: /Search bookmarks/i })
    fireEvent.change(input, { target: { value: 'test@example.com' } })
    vi.advanceTimersByTime(300)

    await waitFor(() => {
      expect(defaultProps.onSearch).toHaveBeenCalledWith('test@example.com')
    })
  })

  it('handles unicode characters in search query', async () => {
    render(<SearchBar {...defaultProps} />)

    const input = screen.getByRole('textbox', { name: /Search bookmarks/i })
    fireEvent.change(input, { target: { value: '日本語テスト' } })
    vi.advanceTimersByTime(300)

    await waitFor(() => {
      expect(defaultProps.onSearch).toHaveBeenCalledWith('日本語テスト')
    })
  })

  it('debounce delay is 300ms', async () => {
    render(<SearchBar {...defaultProps} />)

    const input = screen.getByRole('textbox', { name: /Search bookmarks/i })
    fireEvent.change(input, { target: { value: 'test' } })

    // Advance by 299ms - should not trigger
    vi.advanceTimersByTime(299)
    expect(defaultProps.onSearch).not.toHaveBeenCalled()

    // Advance 1 more ms (total 300ms) - should trigger
    vi.advanceTimersByTime(1)

    await waitFor(() => {
      expect(defaultProps.onSearch).toHaveBeenCalledWith('test')
    })
  })

  // Snapshot test
  it('matches snapshot with empty query', () => {
    const { container } = render(<SearchBar {...defaultProps} />)
    expect(container).toMatchSnapshot()
  })

  it('matches snapshot with query value', () => {
    const { container } = render(<SearchBar {...defaultProps} />)

    const input = screen.getByRole('textbox')
    fireEvent.change(input, { target: { value: 'test query' } })

    expect(container).toMatchSnapshot()
  })
})
