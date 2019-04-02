export const enum AlertType {
  DANGER, WARNING, INFO, SUCCESS
}

export const dismissInterval = 10 * 1000;

export const CommonRoutes = {
  SIGN_IN: '/sign-in',
  F2O_ROOT: '/',
  F2O_DEFAULT: '/fit2openshift'
};

export const httpStatusCode = {
  'Unauthorized': 401,
  'Forbidden': 403
};
