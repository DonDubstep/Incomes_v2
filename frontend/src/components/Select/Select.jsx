import styles from './Select.module.css';

function Select({ value, onChange, options, icon, ariaLabel }) {
  return (
    <label className={styles.wrapper}>
      <select
        className={`${styles.select} ${icon ? styles.withIcon : ''}`}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        aria-label={ariaLabel}
      >
        {options.map((o) => (
          <option key={o.value} value={o.value}>{o.label}</option>
        ))}
      </select>
      {icon && <span className={styles.icon}>{icon}</span>}
      <svg className={styles.chevron} width="12" height="12" viewBox="0 0 12 12" aria-hidden="true">
        <path d="M2 4l4 4 4-4" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" />
      </svg>
    </label>
  );
}

export default Select;
