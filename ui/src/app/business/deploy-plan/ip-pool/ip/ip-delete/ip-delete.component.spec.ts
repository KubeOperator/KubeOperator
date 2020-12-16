import { ComponentFixture, TestBed } from '@angular/core/testing';

import { IpDeleteComponent } from './ip-delete.component';

describe('IpDeleteComponent', () => {
  let component: IpDeleteComponent;
  let fixture: ComponentFixture<IpDeleteComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ IpDeleteComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(IpDeleteComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
