import {Injectable} from '@angular/core';
import {Package} from '../../package/package';
import {PackageService} from '../../package/package.service';
import {Cluster} from '../../cluster/cluster';
import {HttpClient} from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class AppsService {

  constructor(private packageService: PackageService, private http: HttpClient) {

  }


}
