// src/components/Header/Header.jsx
import styles from './Header.module.css';

function Header() {
  return (
    <header className={styles.header}>
      <h1 className={styles.title}>Доходы</h1>
      
      <div className={styles.controls}>
        <div className={styles.datePicker}>
          Текущая дата
        </div>
        <button className={styles.addBtn}>
          <span className={styles.plus}>+</span> Добавить цель
        </button>
      </div>
    </header>
  );
}

export default Header;