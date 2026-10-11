(defun eforth--inner (lang list)
  "Internal function for the `eforth' macro.
LANG are the bindings for the forth.
LIST is full of the forth stack (RPN)."
  (let ((s))
    (cl-flet ((fn (fn) (if (functionp fn)
                           fn
                         (symbol-function fn))))
      (dolist (i list)
        (if (assoc i lang)
            (let ((fn (cadr (assoc i lang))))
              (if s
                  (progn (apply (fn fn) s) (setf s nil))
                (funcall (fn fn))))
          (push i s))))))

(defmacro eforth (lang list)
  "Define and execute a forth language.
LANG are the let-bindings for the FORTH's operators.
LIST is the forth's stack."
  `(eforth--inner ,lang ',list))