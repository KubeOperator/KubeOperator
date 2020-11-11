import {Injectable} from '@angular/core';
import {ActivatedRouteSnapshot, Resolve, RouterStateSnapshot} from '@angular/router';
import {BusinessLicenseService} from './business-license.service';
import {LicenseService} from '../../business/setting/license/license.service';
import {Observable} from 'rxjs';

@Injectable({
    providedIn: 'root'
})
export class BusinessResolverService implements Resolve<boolean> {

    constructor(private licenseService: LicenseService, private businessLicenseService: BusinessLicenseService) {
    }

    resolve(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<boolean> | Promise<boolean> | boolean {
        return this.initLicense().then(data => {
            return data;
        });
    }

    async initLicense() {
        try {
            const data = await this.getLicense();
            if (data && data.status === 'valid') {
                this.businessLicenseService.update(true);
                return true;
            } else {
                return false;
            }
        } catch (e) {
            return false;
        }

    }

    getLicense() {
        return this.licenseService.get().toPromise();
    }
}
