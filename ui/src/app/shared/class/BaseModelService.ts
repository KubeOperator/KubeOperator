import {BaseModel, BaseRequest} from './BaseModel';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {Page} from './Page';
import {Batch} from './Batch';
import { Md5 } from 'ts-md5';

export abstract class BaseModelService<T extends BaseModel> {

    baseUrl = '';
    variable = new Map<string, string>();

    protected constructor(protected http: HttpClient) {
    }

    list(projectName?: string): Observable<Page<T>> {
        const url = this.urlHandler();
        const options = {};
        if (projectName) {
            options['headers'] = {
                project: encodeURI(projectName)
            };
        }
        return this.http.get<Page<T>>(url, options);
    }

    page(page, size, projectName?: string): Observable<Page<T>> {
        const url = this.urlHandler();
        const options = {};
        if (projectName) {
            options['headers'] = {
                project: encodeURI(projectName)
            };
        }
        const pageUrl = `${url}?pageNum=${page}&pageSize=${size}`;
        return this.http.get<Page<T>>(pageUrl, options);
    }

    get(name: string, projectName?: string): Observable<T> {
        const url = this.urlHandler();
        const options = {};
        if (projectName) {
            options['headers'] = {
                project: encodeURI(projectName)
            };
        }
        const itemUrl = `${url}/${name}`;
        return this.http.get<T>(itemUrl, options);
    }

    create(item: BaseRequest, projectName?: string): Observable<T> {
        const url = this.urlHandler();
        const options = {};
        if (projectName) {
            options['headers'] = {
                project: encodeURI(projectName),
                "X-CSRF-TOKEN": this.getCsrf()
            };
        } else {
            options['headers'] = {
                "X-CSRF-TOKEN": this.getCsrf()
            };
        }
        return this.http.post<T>(url, item, options);
    }

    update(name: string, item: BaseRequest, projectName?: string): Observable<T> {
        const url = this.urlHandler();
        const itemUrl = `${url}/${name}/`;
        const options = {};
        if (projectName) {
            options['headers'] = {
                project: encodeURI(projectName)
            };
        }
        return this.http.patch<T>(itemUrl, item, options);
    }

    delete(name: string, projectName?: string): Observable<any> {
        const url = this.urlHandler();
        const options = {};
        if (projectName) {
            options['headers'] = {
                project: encodeURI(projectName)
            };
        }
        const itemUrl = `${url}/${name}/`;
        return this.http.delete<any>(itemUrl, options);
    }


    batch(method: string, items: T[], projectName?: string): Observable<any> {
        const url = this.urlHandler();
        const batchUrl = `${url}/batch/`;
        const options = {};
        if (projectName) {
            options['headers'] = {
                project: encodeURI(projectName)
            };
        }
        const b = new Batch<T>(method, items);
        return this.http.post(batchUrl, b, options);
    }

    private urlHandler() {
        let url = this.baseUrl;
        this.variable.forEach(((k, v) => {
            if (url.indexOf(`{${k}}`) !== -1) {
                url = url.replace(`{${k}}`, v);
            }
        }));
        return url;
    }

    getCsrf(): String {
        var offset = new Date().getTimezoneOffset()
        var thisTime = new Date()
        thisTime.setMinutes(thisTime.getMinutes() + offset)
        let formateDay = (day) => {
          return String(day).replace(/(^\d{1}$)/,'0$1')
        }
        var kk = formateDay(thisTime.getMonth() + 1) + "-" + formateDay(thisTime.getDate()) + " " + formateDay(thisTime.getHours()) + ":" + formateDay(thisTime.getMinutes())
        return (Md5.hashStr("kubeoperator" + kk))
    }
}

