import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ZoneUpdateComponent } from './zone-update.component';

describe('ZoneUpdateComponent', () => {
  let component: ZoneUpdateComponent;
  let fixture: ComponentFixture<ZoneUpdateComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ZoneUpdateComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ZoneUpdateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
