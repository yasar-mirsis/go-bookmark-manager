import { ReactNode } from 'react'
import styles from './Layout.module.css'

interface LayoutProps {
  children: ReactNode
  title?: string
  onAddBookmark?: () => void
  showAddButton?: boolean
}

function Layout({
  children,
  title = 'Bookmark Manager',
  onAddBookmark,
  showAddButton = true,
}: LayoutProps) {
  return (
    <div className={styles.layout}>
      <header className={styles.header}>
        <div className={styles.headerContent}>
          <h1 className={styles.title}>{title}</h1>
          {showAddButton && onAddBookmark && (
            <button
              className={styles.addBookmarkBtn}
              onClick={onAddBookmark}
              type="button"
            >
              <span className={styles.addIcon}>+</span>
              Add Bookmark
            </button>
          )}
        </div>
      </header>

      <main className={styles.main}>{children}</main>

      <footer className={styles.footer}>
        <div className={styles.footerContent}>
          <p className={styles.footerText}>
            Bookmark Manager - Organize your web resources
          </p>
        </div>
      </footer>
    </div>
  )
}

export default Layout
