// src/components/Header/Header.jsx
import styles from './Header.module.css';

function Header({ title = 'Доходы', children }) {
  return (
    <header className={styles.header}>
      <h1 className={styles.title}>{title}</h1>
      
      <div className={styles.controls}>
        {children ?? (
          <>
            <div className={styles.datePicker}>
              Текущая дата
            </div>
            <button className={styles.addBtn}>
              <span className={styles.plus}>+</span> Добавить цель
            </button>
          </>
        )}
      </div>
    </header>
  );
}

export default Header;