import {Injectable} from '@angular/core';
import {Package} from './package';

@Injectable({
  providedIn: 'root'
})
export class PackageLogoService {

  constructor() {
  }

  getLogo(resource: string): string {
    let logo = null;
    const path = 'assets/images';
    switch (resource) {
      case 'kubernetes':
        logo = path + '/logo-k8s.png';
        break;
      case 'okd':
        logo = path + '/logo-okd.png';
        break;
      default:
        logo = 'assets/images/favicon.ico';
        break;
    }
    return logo;
  }
}
