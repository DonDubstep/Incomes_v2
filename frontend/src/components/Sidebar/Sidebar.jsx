// src/components/Sidebar/Sidebar.jsx
import styles from './Sidebar.module.css';

const menuItems = [
  { id: 1, name: 'Доходы', icon: '💰', page: 'incomes' },
  { id: 2, name: 'Статистика', icon: '📊', page: 'stats' },
  { id: 3, name: 'Путешествия', icon: '✈️' },
  { id: 4, name: 'Цели и планы', icon: '🎯' },
];

function Sidebar({ activePage, onNavigate }) {
  return (
    <aside className={styles.sidebar}>
      <div className={styles.logo}>Incomes</div>
      <nav className={styles.nav}>
        <ul>
          {menuItems.map((item) => (
            <li
              key={item.id}
              role="button"
              tabIndex={0}
              className={`${styles.navItem} ${item.page === activePage ? styles.active : ''}`}
              onClick={() => item.page && onNavigate(item.page)}
              onKeyDown={(e) => e.key === 'Enter' && item.page && onNavigate(item.page)}
            >
              <span className={styles.icon}>{item.icon}</span>
              <span className={styles.name}>{item.name}</span>
            </li>
          ))}
        </ul>
      </nav>
      <div className={styles.settings}>
        <span className={styles.icon}>⚙️</span>
        <span className={styles.name}>Настройки</span>
      </div>
    </aside>
  );
}

export default Sidebar;
