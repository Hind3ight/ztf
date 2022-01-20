import request from '@/utils/request';
import { Execution } from './data.d';
import { QueryParams } from '@/types/data.d';
import {SetWidth} from "@/utils/dom";

const apiPath = 'exec';

export async function query(params?: QueryParams): Promise<any> {
    return request({
        url: `/${apiPath}`,
        method: 'get',
        params,
    });
}

export async function create(params: Omit<Execution, 'id'>): Promise<any> {
    return request({
        url: `/${apiPath}`,
        method: 'POST',
        data: params,
    });
}

export async function update(id: number, params: Omit<Execution, 'id'>): Promise<any> {
    return request({
        url: `/${apiPath}/${id}`,
        method: 'PUT',
        data: params,
    });
}

export async function remove(id: number): Promise<any> {
    return request({
        url: `/${apiPath}/${id}`,
        method: 'delete',
    });
}

export async function detail(id: number): Promise<any> {
    return request({url: `/executions/${id}`});
}

export function genExecInfo(jsn: any, i: number): string {
    let msg = jsn.msg.replace(/^"+/,'').replace(/"+$/,'')
    msg = SetWidth(i + '. ', 40) + `<span>${msg}</span>`;

    let sty = ''
    if (jsn.category === 'exec') {
        sty = 'color: #009688;'
    } else if (jsn.category === 'output') {
        // sty = 'font-style: italic;'
    }

    msg = `<div style="${sty}"> ${msg} </div>`

    return msg
}