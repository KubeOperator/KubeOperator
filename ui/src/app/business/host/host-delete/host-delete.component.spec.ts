import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { HostDeleteComponent } from './host-delete.component';

describe('HostDeleteComponent', () => {
  let component: HostDeleteComponent;
  let fixture: ComponentFixture<HostDeleteComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ HostDeleteComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HostDeleteComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
