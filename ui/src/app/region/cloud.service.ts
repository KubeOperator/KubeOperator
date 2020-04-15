import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {CloudZone} from './cloud';
import {ModelMeta, Region} from './region';

@Injectable({
  providedIn: 'root'
})
export class CloudService {
  regionUrl = '/api/v1/cloud/region/';
  zoneUrl = '/api/v1/cloud/{region}/zone/';
  flavorUrl = '/api/v1/cloud/{region}/flavor/';
  templateUrl = '/api/v1/cloud/{region}/template/';

  constructor(private http: HttpClient) {
  }

  listRegion(region: Region): Observable<string[]> {
    return this.http.post<string[]>(this.regionUrl, region);
  }

  listZone(region: string): Observable<CloudZone[]> {
    return this.http.get<CloudZone[]>(this.zoneUrl.replace('{region}', region));
  }

  listFlavor(region: string): Observable<ModelMeta[]> {
    return this.http.get<ModelMeta[]>(this.flavorUrl.replace('{region}', region));
  }

  listTemplates(region: string): Observable<any> {
    return this.http.get<any>(this.templateUrl.replace('{region}', region));
  }
}
