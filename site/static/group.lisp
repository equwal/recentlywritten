(in-package #:cl-user)

(defpackage #:group
  (:use #:cl #:alexandria)
  (:shadowing-import-from :alexandria :once-only)
  (:export :group :test))

(in-package #:group)

(defun group-aux (n source)
  (labels ((inner (down up)
             (if (null down)
                 (nreverse up)
                 (inner (nthcdr n down)
                        (push (subseq down 0 n) up)))))
    (let ((l (length source)))
      (if (= 0 (mod (length source) n))
          (inner source nil)
          (error "Wrong length ~S of source list for groups of size ~S"
                 l n)))))

(defmacro group (x y)
  (if (numberp x)
      (alexandria:once-only (x y)
        `(group-aux ,y ,x))
      `(group-aux ,x ,y)))

(defun test ()
  (let ((res (list (equal (group 2 '(1 2)) '((1 2)))
                   (equal (group '(1 2) 2) '((1 2)))
                   (GROUP '(1 2) (+ 1 1)))))
    (if (every #'identity res)
        (format t "Group does pass tests. Results: ~S" res)
        (error "Group does not pass tests. Results: ~S" res))))
