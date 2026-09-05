// src/components/Sidebar/Sidebar.jsx
import styles from './Sidebar.module.css';

function Sidebar() {
  const menuItems = [
    { id: 1, name: 'Доходы', icon: '💰', active: true },
    { id: 2, name: 'Статистика', icon: '📊' },
    { id: 3, name: 'Путешествия', icon: '✈️' },
    { id: 4, name: 'Цели и планы', icon: '🎯' },
  ];

  return (
    <aside className={styles.sidebar}>
      <div className={styles.logo}>Incomes</div>
      <nav className={styles.nav}>
        <ul>
          {menuItems.map((item) => (
            <li 
              key={item.id} 
              className={`${styles.navItem} ${item.active ? styles.active : ''}`}
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