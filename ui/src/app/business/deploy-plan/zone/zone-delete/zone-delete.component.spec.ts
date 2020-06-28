import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ZoneDeleteComponent } from './zone-delete.component';

describe('ZoneDeleteComponent', () => {
  let component: ZoneDeleteComponent;
  let fixture: ComponentFixture<ZoneDeleteComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ZoneDeleteComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ZoneDeleteComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
