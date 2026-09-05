// src/components/Header/Header.jsx
import styles from './Header.module.css';

function Header() {
  return (
    <header className={styles.header}>
      <h1 className={styles.title}>Доходы</h1>
      
      <div className={styles.controls}>
        <div className={styles.datePicker}>
          📅 1 мая — 31 мая 2024
        </div>
        
        <div className={styles.actions}>
          <div className={styles.notification}>
            🔔<span className={styles.badge}></span>
          </div>
          <img 
            src="https://via.placeholder.com/40" // Или ссылка на реальное фото
            alt="User profile" 
            className={styles.avatar}
          />
        </div>
      </div>
    </header>
  );
}

export default Header;